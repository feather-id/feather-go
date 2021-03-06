package feather

import (
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	featherIssuer = "feather.id"
)

// SessionStatus represents the status of a session.
type SessionStatus string

const (
	// The session is currently active.
	SessionStatusActive = "active"

	// The session has expired.
	SessionStatusExpired = "expired"

	// The session has been revoked.
	SessionStatusRevoked = "revoked"
)

// Session is the Feather session object.
// https://feather.id/docs/reference/api#sessionObject
type Session struct {
	ID        string        `json:"id"`
	Object    string        `json:"object"`
	Status    SessionStatus `json:"status"`
	Token     *string       `json:"token"`
	UserID    string        `json:"user_id"`
	CreatedAt time.Time     `json:"created_at"`
	RevokedAt *time.Time    `json:"revoked_at"`
}

// SessionList is a list of Feather session objects.
// https://feather.id/docs/reference/api#pagination
type SessionList struct {
	ListMeta
	Data []*Session `json:"data"`
}

// Sessions provides an interface for accessing Feather API session objects.
// https://feather.id/docs/reference/api#sessions
type Sessions interface {
	Create(params SessionsCreateParams) (*Session, error)
	List(params SessionsListParams) (*SessionList, error)
	Retrieve(id string) (*Session, error)
	Revoke(id string, params SessionsRevokeParams) (*Session, error)
	Upgrade(id string, params SessionsUpgradeParams) (*Session, error)
	Validate(params SessionsValidateParams) (*Session, error)
}

type sessions struct {
	gateway          gateway
	cachedPublicKeys map[string]*rsa.PublicKey
}

func newSessionsResource(g gateway) sessions {
	return sessions{
		gateway:          g,
		cachedPublicKeys: map[string]*rsa.PublicKey{},
	}
}

// Create a new session.
// https://feather.id/docs/reference/api#createSession
func (s sessions) Create(params SessionsCreateParams) (*Session, error) {
	var session Session
	if err := s.gateway.sendRequest(http.MethodPost, pathSessions, params, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// SessionsCreateParams ...
type SessionsCreateParams struct {
	CredentialToken *string `json:"credential_token"`
}

// List a user's sessions.
// https://feather.id/docs/reference/api#listSessions
func (s sessions) List(params SessionsListParams) (*SessionList, error) {
	var sessionList SessionList
	if err := s.gateway.sendRequest(http.MethodGet, pathSessions, params, &sessionList); err != nil {
		return nil, err
	}
	return &sessionList, nil
}

// SessionsListParams ...
type SessionsListParams struct {
	ListParams
	UserID *string `json:"user_id"`
}

// Retrieve a session.
// https://feather.id/docs/reference/api#retrieveSession
func (s sessions) Retrieve(id string) (*Session, error) {
	var session Session
	path := strings.Join([]string{pathSessions, id}, "/")
	if err := s.gateway.sendRequest(http.MethodGet, path, nil, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// Revoke a session.
// https://feather.id/docs/reference/api#revokeSession
func (s sessions) Revoke(id string, params SessionsRevokeParams) (*Session, error) {
	var session Session
	path := strings.Join([]string{pathSessions, id, "revoke"}, "/")
	if err := s.gateway.sendRequest(http.MethodPost, path, params, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// SessionsRevokeParams ...
type SessionsRevokeParams struct {
	SessionToken *string `json:"session_token"`
}

// Upgrade a session.
// https://feather.id/docs/reference/api#upgradeSession
func (s sessions) Upgrade(id string, params SessionsUpgradeParams) (*Session, error) {
	var session Session
	path := strings.Join([]string{pathSessions, id, "upgrade"}, "/")
	if err := s.gateway.sendRequest(http.MethodPost, path, params, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// SessionsUpgradeParams ...
type SessionsUpgradeParams struct {
	CredentialToken *string `json:"credential_token"`
}

// Validate a session.
// https://feather.id/docs/reference/api#validateSession
func (s sessions) Validate(params SessionsValidateParams) (*Session, error) {
	if params.SessionToken == nil {
		return nil, Error{
			Type:    ErrorTypeValidation,
			Code:    ErrorCodeSessionTokenInvalid,
			Message: "No session tokens were not provided for validation",
		}
	}

	session, err := s.parseSessionToken(*params.SessionToken)
	if err != nil {
		ferr, _ := err.(Error)
		if ferr.Code == ErrorCodeSessionTokenExpired {
			// TODO send the session token to the API
			path := strings.Join([]string{pathSessions, session.ID, "validate"}, "/")
			if err := s.gateway.sendRequest(http.MethodPost, path, params, session); err != nil {
				return nil, err
			}
		} else {
			return nil, ferr
		}
	}

	return session, nil
}

// SessionsValidateParams ...
type SessionsValidateParams struct {
	SessionToken *string `json:"session_token"`
}

func (s *sessions) parseSessionToken(tokenStr string) (*Session, error) {
	invalidTokenError := Error{
		Object:  "error",
		Type:    ErrorTypeValidation,
		Code:    ErrorCodeSessionTokenInvalid,
		Message: "The session token is invalid",
	}

	// Parse the string for a token
	parser := jwt.Parser{
		ValidMethods:         []string{jwt.SigningMethodRS256.Name},
		SkipClaimsValidation: true,
	}

	// Parse the token
	token, err := parser.Parse(tokenStr, s.getValidationKey)
	if err != nil {
		return nil, invalidTokenError
	}

	// Validate the basic session token claims
	claims := token.Claims.(jwt.MapClaims)
	if claims["iss"] != featherIssuer {
		return nil, invalidTokenError
	}
	subject, ok := claims["sub"].(string)
	if !ok || !strings.HasPrefix(subject, "USR_") {
		return nil, invalidTokenError
	}
	audience, ok := claims["aud"].(string)
	if !ok || !strings.HasPrefix(audience, "PRJ_") {
		return nil, invalidTokenError
	}
	sessionID, ok := claims["ses"].(string)
	if !ok || !strings.HasPrefix(sessionID, "SES_") {
		return nil, invalidTokenError
	}
	cat, ok := claims["cat"].(float64)
	if !ok {
		return nil, invalidTokenError
	}
	createdAt := time.Unix(int64(cat), 0).UTC()

	// Generate a session object from the token
	session := Session{
		ID:        sessionID,
		Object:    "session",
		Status:    SessionStatusActive,
		Token:     &tokenStr,
		UserID:    subject,
		CreatedAt: createdAt,
		RevokedAt: nil,
	}

	// Check if the token is expired
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, invalidTokenError
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		return &session, Error{
			Type:    ErrorTypeValidation,
			Code:    ErrorCodeSessionTokenExpired,
			Message: "The provided session token is expired",
		}
	}

	return &session, nil
}

func (s *sessions) getValidationKey(token *jwt.Token) (interface{}, error) {
	keyID, ok := token.Header["kid"].(string)
	if !ok || keyID == "" {
		return nil, errors.New("Header 'kid' not found")
	}
	return s.getPublicKey(keyID)
}

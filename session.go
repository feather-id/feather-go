package feather

import (
	"net/http"
	"strings"
	"time"
)

// TODO Revoke

// Session is the Feather session object.
// https://feather.id/docs/reference/api#sessionObject
type Session struct {
	ID        string     `json:"id"`
	Object    string     `json:"object"`
	Type      string     `json:"type"`   // TODO make enum
	Status    string     `json:"status"` // TODO make enum
	Token     *string    `json:"token"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
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
	Upgrade(id string, params SessionsUpgradeParams) (*Session, error)
	Validate(params SessionsValidateParams) (*Session, error)
}

type sessions struct {
	gateway gateway
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
	panic("not implemented")
}

// SessionsValidateParams ...
type SessionsValidateParams struct {
	SessionToken *string `json:"session_token"`
}

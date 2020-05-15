package feather

import (
	"time"
)

// Session is the Feather session object.
// https://feather.id/docs/reference/api#sessionObject
type Session struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`   // TODO make enum
	Status    string     `json:"status"` // TODO make enum
	Token     *string    `json:"token"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

// Sessions provides an interface for accessing Feather API session objects.
// https://feather.id/docs/reference/api#sessions
type Sessions interface {
	Create(params SessionsCreateParams) (*Session, error)
	List(param SessionsListParams) // TODO lists
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
	panic("not implemented")
}

// SessionsCreateParams ...
type SessionsCreateParams struct {
	CredentialToken string `json:"credential_token"`
}

// List a user's sessions.
// https://feather.id/docs/reference/api#listSessions
func (s sessions) List(param SessionsListParams) {
	panic("not implemented")
}

// SessionsListParams ...
type SessionsListParams struct {
	UserID        string `json:"user_id"`
	Limit         int    `json:"limit"`
	StartingAfter string `json:"starting_after"`
	EndingBefore  string `json:"ending_before"`
}

// Retrieve a session.
// https://feather.id/docs/reference/api#retrieveSession
func (s sessions) Retrieve(id string) (*Session, error) {
	panic("not implemented")
}

// Upgrade a session.
// https://feather.id/docs/reference/api#upgradeSession
func (s sessions) Upgrade(id string, params SessionsUpgradeParams) (*Session, error) {
	panic("not implemented")
}

// SessionsUpgradeParams ...
type SessionsUpgradeParams struct {
	CredentialToken string `json:"credential_token"`
}

// Validate a session.
// https://feather.id/docs/reference/api#validateSession
func (s sessions) Validate(params SessionsValidateParams) (*Session, error) {
	panic("not implemented")
}

// SessionsValidateParams ...
type SessionsValidateParams struct {
	SessionToken string `json:"session_token"`
}

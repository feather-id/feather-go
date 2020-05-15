package feather

import (
	"feather/IDPLAT/src/common/fid"
	"time"
)

// Sessions ...
var Sessions = sessionsResource{}

// Sessions resource

type sessionsResource struct {
}

func (s sessionsResource) Create(params SessionsCreateParams) (*Session, error) {
	panic("not implemented")
}

func (s sessionsResource) List(param SessionsListParams) {
	panic("not implemented")
}

func (s sessionsResource) Retrieve(id string) (*Session, error) {
	panic("not implemented")
}

func (s sessionsResource) Upgrade(id string, params SessionsUpgradeParams) (*Session, error) {
	panic("not implemented")
}

func (s sessionsResource) Validate(params SessionsValidateParams) (*Session, error) {
	panic("not implemented")
}

// Session is the Feather session object
// https://feather.id/docs/reference/api#sessions
type Session struct {
	ID        fid.FID    `json:"id"`
	Object    string     `json:"object"`
	Type      string     `json:"type"`   // TODO make enum
	Status    string     `json:"status"` // TODO make enum
	Token     *string    `json:"token"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

// SessionsCreateParams ...
type SessionsCreateParams struct {
	CredentialToken string `json:"credential_token"`
}

// SessionsListParams ...
type SessionsListParams struct {
	UserID        string `json:"user_id"`
	Limit         int    `json:"limit"`
	StartingAfter string `json:"starting_after"`
	EndingBefore  string `json:"ending_before"`
}

// SessionsUpgradeParams ...
type SessionsUpgradeParams struct {
	CredentialToken string `json:"credential_token"`
}

// SessionsValidateParams ...
type SessionsValidateParams struct {
	SessionToken string `json:"session_token"`
}

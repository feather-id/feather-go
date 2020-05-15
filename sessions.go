package feather

import (
	"time"
)

// Sessions resource interface ...
type Sessions interface {
	Create(params SessionsCreateParams) (*Session, error)
	List(param SessionsListParams)
	Retrieve(id string) (*Session, error)
	Upgrade(id string, params SessionsUpgradeParams) (*Session, error)
	Validate(params SessionsValidateParams) (*Session, error)
}

type sessions struct {
	gateway gateway
}

func (s sessions) Create(params SessionsCreateParams) (*Session, error) {
	panic("not implemented")
}

func (s sessions) List(param SessionsListParams) {
	panic("not implemented")
}

func (s sessions) Retrieve(id string) (*Session, error) {
	panic("not implemented")
}

func (s sessions) Upgrade(id string, params SessionsUpgradeParams) (*Session, error) {
	panic("not implemented")
}

func (s sessions) Validate(params SessionsValidateParams) (*Session, error) {
	panic("not implemented")
}

// Session is the Feather session object
// https://feather.id/docs/reference/api#sessions
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

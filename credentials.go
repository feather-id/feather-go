package feather

import (
	"time"
)

// Credentials ...
var Credentials = credentialsResource{}

type credentialsResource struct {
}

func (c credentialsResource) Create(params CredentialsCreateParams) (*Credential, error) {
	panic("not implemented")
}

func (c credentialsResource) Update(id string, params CredentialsUpdateParams) (*Credential, error) {
	panic("not implemented")
}

// Credential is a Feather credential object
type Credential struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"` // TODO make enum
	Token     *string   `json:"token"`
	Type      string    `json:"type"` // TODO make enum
}

// CredentialsCreateParams ...
type CredentialsCreateParams struct {
	Type     string `json:"type"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CredentialsUpdateParams ...
type CredentialsUpdateParams struct {
	OneTimeCode string `json:"one_time_code"`
}

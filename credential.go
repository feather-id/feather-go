package feather

import (
	"net/http"
	"strings"
	"time"
)

// Credential is a Feather credential object
// https://feather.id/docs/reference/api#credentialObject
type Credential struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"` // TODO make enum
	Token     *string   `json:"token"`
	Type      string    `json:"type"` // TODO make enum
}

// Credentials provides a set of operations to access a Feather API credential resource.
// https://feather.id/docs/reference/api#credentials
type Credentials interface {
	Create(params CredentialsCreateParams) (*Credential, error)
	Update(id string, params CredentialsUpdateParams) (*Credential, error)
}

type credentials struct {
	gateway gateway
}

// Create creates a new Credential object
// https://feather.id/docs/reference/api#createCredential
func (c credentials) Create(params CredentialsCreateParams) (*Credential, error) {
	var credential Credential
	if err := c.gateway.sendRequest(http.MethodPost, resourcePathCredentials, params, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

// CredentialsCreateParams ...
type CredentialsCreateParams struct {
	Type     string  `json:"type"` // TODO make enum
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// Update updates a Credential object
// https://feather.id/docs/reference/api#updateCredential
func (c credentials) Update(id string, params CredentialsUpdateParams) (*Credential, error) {
	var credential Credential
	path := strings.Join([]string{resourcePathCredentials, id}, "/")
	if err := c.gateway.sendRequest(http.MethodPost, path, params, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

// CredentialsUpdateParams ...
type CredentialsUpdateParams struct {
	OneTimeCode *string `json:"one_time_code"`
}

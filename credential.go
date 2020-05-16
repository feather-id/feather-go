package feather

import (
	"net/http"
	"strings"
	"time"
)

// CredentialStatus represents the status of a credential.
type CredentialStatus string

const (
	// The provided authentication information was valid.
	CredentialStatusValid = "valid"

	// The provided authentication information was invalid.
	CredentialStatusInvalid = "invalid"

	// A one-time-code has been sent to the user and must be returned
	// to verify the provided authentication information.
	CredentialStatusRequiresOneTimeCode = "requires_one_time_code"
)

// CredentialType represents the type of the provided authentication information.
type CredentialType string

const (
	// Only an email address was provided.
	CredentialTypeEmail = "email"

	// An email address and password were provided.
	CredentialTypeEmailPassword = "email|password"

	// A username and password were provided.
	CredentialTypeUsernamePassword = "username|password"
)

// Credential is a Feather credential object.
// https://feather.id/docs/reference/api#credentialObject
type Credential struct {
	ID        string           `json:"id"`
	Object    string           `json:"object"`
	CreatedAt time.Time        `json:"created_at"`
	ExpiresAt time.Time        `json:"expires_at"`
	Status    CredentialStatus `json:"status"`
	Token     *string          `json:"token"`
	Type      CredentialType   `json:"type"` // TODO make enum
}

// Credentials provides an interface for accessing Feather API credential objects.
// https://feather.id/docs/reference/api#credentials
type Credentials interface {
	Create(params CredentialsCreateParams) (*Credential, error)
	Update(id string, params CredentialsUpdateParams) (*Credential, error)
}

type credentials struct {
	gateway gateway
}

// Create a new credential.
// https://feather.id/docs/reference/api#createCredential
func (c credentials) Create(params CredentialsCreateParams) (*Credential, error) {
	var credential Credential
	if err := c.gateway.sendRequest(http.MethodPost, pathCredentials, params, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

// CredentialsCreateParams ...
type CredentialsCreateParams struct {
	Type     CredentialType `json:"type"`
	Email    *string        `json:"email"`
	Username *string        `json:"username"`
	Password *string        `json:"password"`
}

// Update a credential.
// https://feather.id/docs/reference/api#updateCredential
func (c credentials) Update(id string, params CredentialsUpdateParams) (*Credential, error) {
	var credential Credential
	path := strings.Join([]string{pathCredentials, id}, "/")
	if err := c.gateway.sendRequest(http.MethodPost, path, params, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

// CredentialsUpdateParams ...
type CredentialsUpdateParams struct {
	OneTimeCode *string `json:"one_time_code"`
}

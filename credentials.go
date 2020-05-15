package feather

import (
	"net/http"
	"strings"
	"time"
)

// Credentials ...
var Credentials = credentialsResource{}

type credentialsResource struct{}

func (c credentialsResource) Create(params CredentialsCreateParams) (*Credential, error) {
	var credential Credential
	if err := apiGateway.sendRequest(http.MethodPost, resourcePathCredentials, params, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

func (c credentialsResource) Update(id string, params CredentialsUpdateParams) (*Credential, error) {
	var credential Credential
	path := strings.Join([]string{resourcePathCredentials, id}, "/")
	if err := apiGateway.sendRequest(http.MethodPost, path, params, &credential); err != nil {
		return nil, err
	}
	return &credential, nil
}

// Credential is a Feather credential object
// https://feather.id/docs/reference/api#credentials
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
	Type     string  `json:"type"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// CredentialsUpdateParams ...
type CredentialsUpdateParams struct {
	OneTimeCode *string `json:"one_time_code"`
}

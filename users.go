package feather

import (
	"net/http"
	"strings"
	"time"
)

// User is the Feather user object.
// https://feather.id/docs/reference/api#userObject
type User struct {
	ID        string    `json:"id"`
	Email     *string   `json:"email"`
	Username  *string   `json:"username"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Users provides an interface for accessing Feather API user objects.
// https://feather.id/docs/reference/api#users
type Users interface {
	List(params UsersListParams) // TODO lists
	Retrieve(id string) (*User, error)
	Update(id string, params UsersUpdateParams) (*User, error)
}

type users struct {
	gateway gateway
}

// List a project's users.
// https://feather.id/docs/reference/api#listUsers
func (u users) List(params UsersListParams) {
	panic("not implemented")
}

// UsersListParams ...
type UsersListParams struct {
	Limit         int    `json:"limit"`
	StartingAfter string `json:"starting_after"`
	EndingBefore  string `json:"ending_before"`
}

// Retrieve a user.
// https://feather.id/docs/reference/api#retrieveUser
func (u users) Retrieve(id string) (*User, error) {
	var user User
	path := strings.Join([]string{pathUsers, id}, "/")
	if err := u.gateway.sendRequest(http.MethodGet, path, nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Update a user.
// https://feather.id/docs/reference/api#updateUser
func (u users) Update(id string, params UsersUpdateParams) (*User, error) {
	var user User
	path := strings.Join([]string{pathUsers, id}, "/")
	if err := u.gateway.sendRequest(http.MethodPost, path, params, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UsersUpdateParams ...
type UsersUpdateParams struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Metadata string  `json:"metadata"`
}

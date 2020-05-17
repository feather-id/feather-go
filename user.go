package feather

import (
	"net/http"
	"strings"
	"time"
)

// User is the Feather user object.
// https://feather.id/docs/reference/api#userObject
type User struct {
	ID        string            `json:"id"`
	Object    string            `json:"object"`
	Email     *string           `json:"email"`
	Username  *string           `json:"username"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// UserList is a list of Feather user objects.
// https://feather.id/docs/reference/api#pagination
type UserList struct {
	ListMeta
	Data []*User `json:"data"`
}

// Users provides an interface for accessing Feather API user objects.
// https://feather.id/docs/reference/api#users
type Users interface {
	List(params UsersListParams) (*UserList, error)
	Retrieve(id string) (*User, error)
	Update(id string, params UsersUpdateParams) (*User, error)
}

type users struct {
	gateway gateway
}

func newUsersResource(g gateway) users {
	return users{
		gateway: g,
	}
}

// List a project's users.
// https://feather.id/docs/reference/api#listUsers
func (u users) List(params UsersListParams) (*UserList, error) {
	var userList UserList
	if err := u.gateway.sendRequest(http.MethodGet, pathUsers, params, &userList); err != nil {
		return nil, err
	}
	return &userList, nil
}

// UsersListParams ...
type UsersListParams struct {
	ListParams
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
	Email    *string            `json:"email"`
	Username *string            `json:"username"`
	Metadata *map[string]string `json:"metadata"`
}

package feather

import "time"

// Users resesource interface ...
var Users = usersResource{}

type usersResource struct {
}

func (u usersResource) List(params UsersListParams) {
	panic("not implemented")
}

func (u usersResource) Retrieve(id string) (*User, error) {
	panic("not implemented")
}

func (u usersResource) Update(id string, params UsersUpdateParams) (*User, error) {
	panic("not implemented")
}

// User is the Feather user object
// https://feather.id/docs/reference/api#users
type User struct {
	ID        string    `json:"id"`
	Object    string    `json:"object"`
	Email     *string   `json:"email"`
	Username  *string   `json:"username"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersListParams ...
type UsersListParams struct {
	Limit         int    `json:"limit"`
	StartingAfter string `json:"starting_after"`
	EndingBefore  string `json:"ending_before"`
}

// UsersUpdateParams ...
type UsersUpdateParams struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Metadata string  `json:"metadata"`
}

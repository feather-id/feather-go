package feather

import (
	"github.com/feather-id/feather-go/resource"
)

// New initializes a new Feather client
func New(apiKey string, config Config) (Feather, error) {
	return Feather{
		Credentials: resource.NewCredentials(),
		PublicKeys:  resource.NewPublicKeys(),
		Sessions:    resource.NewSessions(),
		Users:       resource.NewUsers(),
		apiKey:      apiKey,
		config:      config,
	}, nil
}

// Feather ...
type Feather struct {
	// Public
	Credentials resource.Credentials
	PublicKeys  resource.PublicKeys
	Sessions    resource.Sessions
	Users       resource.Users

	// Private
	apiKey string
	config Config
}

// Config allows the caller to configure the Feather SDK setup
type Config struct {
	Host     string
	Port     string
	Protocol string
	BasePath string
}

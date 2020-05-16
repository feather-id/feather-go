// Package feather provides a convenient interface to the Feather API for applications running in a Go server environment.
//
// For more information on Feather API, please check out our docs at https://feather.id/docs.
package feather

import "net/http"

// A Client provides access to the Feather API core resources.
// You should instantiate and use a client to send requests to
// the Feather API.
//
// Example:
//     // Create a Feather client
//     client := feather.New("test_ABC")
//
//     // Create an anonymous session
//     credential, err := client.Sessions.Create(null)
//     if err != nil {
//         // Handle error
//     }
type Client struct {
	Credentials Credentials
	Sessions    Sessions
	Users       Users
}

// A Config provides extra configuration to intialize a Feather client with.
// This is typically only needed in a testing/development environment and should
// not be used in production code.
type Config struct {
	Protocol   *string
	Host       *string
	Port       *string
	BasePath   *string
	HTTPClient *http.Client
}

// New creates a new instance of the Feather client.
// If additional configuration is needed for the client instance,
// use the optional Config parameter to add the extra config.
//
// Example:
//     // Create a Feather client with just an API key
//     client := feather.New("test_ABC")
//
//     // Create a Feather client with additional configuration
//     client := feather.New("test_ABC", &feather.Config{
//         Host: feather.String("localhost"),
//     })
func New(apiKey string, cfgs ...*Config) Client {
	cfg := Config{}
	if len(cfgs) > 0 {
		cfg = *cfgs[0]
	}
	g := gateway{
		apiKey: apiKey,
		config: cfg,
	}
	return Client{
		Credentials: credentials{gateway: g},
		Sessions:    sessions{gateway: g},
		Users:       users{gateway: g},
	}
}

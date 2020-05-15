package feather

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
	Host     string
	Port     string
	Protocol string
	BasePath string
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
//     client := feather.New("test_ABC", feather.Config{ Host: "localhost" })
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

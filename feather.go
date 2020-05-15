package feather

// APIKey the key the API will be called with
var APIKey string

// SetConfig ...
func SetConfig(cfg Config) {
	apiGateway.config = cfg
}

// Config allows the caller to configure the Feather SDK setup
type Config struct {
	Host     string
	Port     string
	Protocol string
	BasePath string
}

const (
	resourcePathCredentials = "/credentials"
	resourcePathSessions    = "/sessions"
	resourcePathUsers       = "/users"
)

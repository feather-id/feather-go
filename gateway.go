package feather

const (
	defaultProtocol string = "https"
	defaultHost     string = "api.feather.id"
	defaultPort     string = "443"
	defaultBasePath string = "/v1"
)

type gateway struct{}

func (g gateway) sendRequest(method string, path string, data interface{}, cfg Config) error {
	panic("Not implemented")
}

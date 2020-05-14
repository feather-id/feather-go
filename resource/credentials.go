package resource

// Credentials resource interface ...
type Credentials interface {
	Create()
	Update()
}

type credentials struct {
}

// NewCredentials ...
func NewCredentials() Credentials {
	return &credentials{}
}

func (c credentials) Create() {
	// TODO
}

func (c credentials) Update() {
	// TODO
}

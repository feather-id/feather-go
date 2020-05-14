package resource

// PublicKeys resource interface ...
type PublicKeys interface {
	Retrieve()
}

type publicKeys struct {
}

// NewPublicKeys ...
func NewPublicKeys() PublicKeys {
	return &publicKeys{}
}

func (s publicKeys) Retrieve() {
	// TODO
}

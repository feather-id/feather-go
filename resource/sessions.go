package resource

// Sessions resource interface ...
type Sessions interface {
	Create()
	List()
	Retrieve()
	Update()
	Validate()
}

type sessions struct {
}

// NewSessions ...
func NewSessions() Sessions {
	return &sessions{}
}

func (s sessions) Create() {
	// TODO
}

func (s sessions) List() {
	// TODO
}

func (s sessions) Retrieve() {
	// TODO
}

func (s sessions) Update() {
	// TODO
}

func (s sessions) Validate() {
	// TODO
}

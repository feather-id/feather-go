package resource

// Users resesource interface ...
type Users interface {
	List()
	Retrieve()
	Update()
}

type users struct {
}

// NewUsers ...
func NewUsers() Users {
	return &users{}
}

func (u users) List() {
	// TODO
}

func (u users) Retrieve() {
	// TODO
}

func (u users) Update() {
	// TODO
}

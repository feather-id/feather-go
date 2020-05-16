package feather

// Error is the Feather error object
// https://feather.id/docs/reference/api#errors
type Error struct {
	Object  string `json:"object"`
	Type    string `json:"type"` // TODO make enum
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

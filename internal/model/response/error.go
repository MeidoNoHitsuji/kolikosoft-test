package response

type Error struct {
	Message string `json:"message"`
}

func NewError(err error) *Error {
	return &Error{Message: err.Error()}
}

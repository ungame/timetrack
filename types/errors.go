package types

type Error struct {
	Err string `json:"error"`
}

func NewError(e error) *Error {
	err := Error{"error"}
	if e != nil {
		err.Err = e.Error()
	}
	return &err
}

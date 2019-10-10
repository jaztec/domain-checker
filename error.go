package checker

import "fmt"

// Error defines a structured error this package will use
type Error struct {
	registrar Registrar
	err       error
}

func (e Error) Error() string {
	return fmt.Sprintf("%T: %s", e.registrar, e.err)
}

// NewError returns an structured error
func NewError(client Registrar, err error) Error {
	return Error{client, err}
}

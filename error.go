package checker

import "fmt"

// Error defines a structured error this package will use
type Error struct {
	client Client
	err    error
}

func (e Error) Error() string {
	return fmt.Sprintf("%T: %s", e.client, e.err)
}

// NewError returns an structured error
func NewError(client Client, err error) Error {
	return Error{client, err}
}

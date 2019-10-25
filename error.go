package checker

import (
	"errors"
	"fmt"
)

// Error defines a structured error this package will use
type Error struct {
	registrar Registrar
	err       error
}

func (e Error) Error() string {
	return fmt.Sprintf("%T: %s", e.registrar, e.err)
}

// Unwrap fills the go 1.13 error interface for chaining
func (e Error) Unwrap() error {
	return e.err
}

// MultipleError holds a set of errors
type MultipleError struct {
	msg  string
	errs []Error
}

func (me *MultipleError) Error() string {
	r := me.msg
	for _, e := range me.errs {
		r += "\n\t- " + e.Error()
	}
	return r
}

// Errors returns all the internal errors contained by this
// multiple error bucket
func (me *MultipleError) Errors() []Error {
	return me.errs
}

// Is implements the interface the errors package can use to
// match the MultipleError to an error that should be tested.
func (me *MultipleError) Is(target error) (matches bool) {
	for _, e := range me.errs {
		matches = errors.Is(e, target)
		if matches == true {
			return
		}
	}
	return
}

// As implements the interface the errors package can use to
// check for errors.As
func (me *MultipleError) As(target interface{}) (matches bool) {
	for _, e := range me.errs {
		matches = errors.As(e, target)
		if matches == true {
			return
		}
	}
	return
}

// Add adds an Error instance to the MultipleErrors error set
func (me *MultipleError) Add(e Error) int {
	me.errs = append(me.errs, e)
	return len(me.errs)
}

// Len gets the size of the internal error set
func (me *MultipleError) Len() int {
	return len(me.errs)
}

// NewError returns an structured error
func NewError(client Registrar, err error) Error {
	return Error{client, err}
}

// NewMultipleError returns a new instance of a multiple error
// object. Any procedures that run in batch should use this
// error to represent problems somewhere down the chain. The
// parameter 'cap' lets you set the capacity of the internal
// slice beforehand to prevent memory allocations during filling.
func NewMultipleError(msg string, cap int) *MultipleError {
	return &MultipleError{
		msg:  msg,
		errs: make([]Error, 0, cap),
	}
}

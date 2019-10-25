package checker

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	msg := "random testing error message"
	expect := "checker.errorRegistrar: " + msg
	orig := errors.New(msg)
	err := NewError(errorRegistrar{}, orig)

	got := fmt.Sprintf("%s", err)
	if got != expect {
		t.Logf("Received '%s' but expected '%s'", got, expect)
		t.Fail()
	}

	if err.Unwrap() != orig {
		t.Logf("Received '%v' but expected '%v'", err.Unwrap(), orig)
		t.Fail()
	}
}

func TestMultipleError(t *testing.T) {
	msg := "multiple test error"
	msg2 := "random testing error message"
	err := NewMultipleError(msg, 1)
	orig := errors.New(msg2)
	i := err.Add(NewError(errorRegistrar{}, orig))
	if i != 1 {
		t.Logf("Expected to have added only one element, received %d", i)
		t.Fail()
	}
	if i != err.Len() {
		t.Logf("Expected to have the same amount of elements as returned by add but got %d", err.Len())
		t.Fail()
	}

	t.Run("check internal error", func(t *testing.T) {
		expect := msg + "\n\t- checker.errorRegistrar: " + msg2
		if got := err.Error(); got != expect {
			t.Logf("Expected '%s' but received '%s'", expect, got)
			t.Fail()
		}
	})

	t.Run("peek into internal errors", func(t *testing.T) {
		errs := err.Errors()
		if n := err.Len(); n != len(errs) {
			t.Logf("Object reports %d internal errors but received %d", n, len(errs))
			t.Fail()
		}
	})

	t.Run("validate the internal Is interface implemetation", func(t *testing.T) {
		if !errors.Is(err, orig) {
			t.Log("The MultipleError instance does not correctly matches the original error")
			t.Fail()
		}
	})

	t.Run("validate the internal As interface implemetation", func(t *testing.T) {
		var err2 Error
		success := errors.As(err, &err2)
		if !success {
			t.Log("The MultipleError instance does not correctly matches the original error")
			t.Fail()
		}
		fmt.Println(err2)
	})
}

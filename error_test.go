package checker

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	t.Run("Test creating error", func(t *testing.T) {
		msg := "random testing error message"
		expect := "checker.errorClient: " + msg
		err := NewError(errorClient{}, errors.New(msg))

		if got := fmt.Sprintf("%s", err); got != expect {
			t.Logf("Received '%s' but expected '%s'", got, expect)
			t.Fail()
		}
	})
}
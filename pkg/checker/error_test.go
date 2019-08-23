package checker

import (
	"errors"
	"fmt"
	"testing"
)

type errClient struct {}

func (errClient) CheckDomain(string) (Status, error) { return Unavailable, nil }

func TestNewError (t *testing.T) {
	t.Run("Test creating error", func (t *testing.T) {
		msg := "random testing error message"
		expect := "checker.errClient: " + msg
		err := NewError(errClient{}, errors.New(msg))

		if got := fmt.Sprintf("%s", err); got != expect {
			t.Logf("Received '%s' but expected '%s'", got, expect)
			t.Fail()
		}
	})
}
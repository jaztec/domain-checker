package checker

import (
	"errors"
	"testing"
)

const errorMessage = "some error message"

type availableClient struct{}
type unavailableClient struct{}
type ownedClient struct{}
type errorClient struct{}

func (availableClient) CheckDomain(string) (Status, error)   { return Available, nil }
func (unavailableClient) CheckDomain(string) (Status, error) { return Unavailable, nil }
func (ownedClient) CheckDomain(string) (Status, error)       { return Owned, nil }
func (errorClient) CheckDomain(string) (Status, error)       { return Unavailable, errors.New(errorMessage) }

type expectedCheckDomainResult struct {
	client Client
	status Status
}

func TestCheckDomain(t *testing.T) {
	t.Run("Test statuses", func(t *testing.T) {
		clients := []Client{
			availableClient{},
			unavailableClient{},
			ownedClient{},
			errorClient{},
		}
		name := "irrelevant"
		expectLen := 3
		expectedResults := []expectedCheckDomainResult{
			{clients[0], Available},
			{clients[1], Unavailable},
			{clients[2], Owned},
		}

		statuses := CheckDomain(name, clients)

		if gotLen := len(statuses); gotLen != expectLen {
			t.Logf("Expected %d result statuses but received %d", expectLen, gotLen)
			t.Fail()
		}

		// test the results in order of client ordering, it should comply
		for i, s := range expectedResults {
			status := statuses[i]
			if status.Status() != s.status {
				t.Logf("Expected %T from status but received %T", s.status, status.Status())
				t.Fail()
			}
			if status.Client() != s.client {
				t.Logf("Expected %T from status but received %T", s.client, status.Client())
				t.Fail()
			}
			if status.Domain() != name {
				t.Logf("Expected %s from status but received %s", name, status.Domain())
				t.Fail()
			}
		}
	})
}

package checker

import (
	"testing"
)

const errorMessage = "some error message"

type expectedCheckDomainResult struct {
	client Client
	status Status
}

var (
	clients = []Client{
		availableClient{},
		unavailableClient{},
		ownedClient{},
		processingClient{},
		errorClient{},
	}
	name = "irrelevant"
)

func TestCheckDomain(t *testing.T) {
	t.Run("Test statuses", func(t *testing.T) {

		expectLen := len(clients) - 1
		expectedResults := []expectedCheckDomainResult{
			{clients[0], Available},
			{clients[1], Unavailable},
			{clients[2], Owned},
			{clients[3], Processing},
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

func TestRegisterDomain(t *testing.T) {
	t.Run("Test registering domains", func(t *testing.T) {
		for range clients {
			if _, err := RegisterDomain(name, clients); err != nil {

			}
		}
	})
}

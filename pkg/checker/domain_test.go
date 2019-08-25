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
	domainClients = []Client{
		availableClient{},
		unavailableClient{},
		ownedClient{},
		processingClient{},
		errorClient{},
	}
	registerClientsSuccess = []Client{
		errorClient{},
		unavailableClient{},
		ownedClient{},
	}
	registerClientsFailure = []Client{
		errorClient{},
		unavailableClient{},
	}
	name = "irrelevant"
)

func TestCheckDomain(t *testing.T) {
	t.Run("Test statuses", func(t *testing.T) {

		expectLen := len(domainClients) - 1
		expectedResults := []expectedCheckDomainResult{
			{domainClients[0], Available},
			{domainClients[1], Unavailable},
			{domainClients[2], Owned},
			{domainClients[3], Processing},
		}

		statuses := CheckDomain(name, domainClients)

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
				t.Logf("Expected %T from client but received %T", s.client, status.Client())
				t.Fail()
			}
			if status.Domain() != name {
				t.Logf("Expected %s from domain but received %s", name, status.Domain())
				t.Fail()
			}
		}
	})
}

func TestRegisterDomain(t *testing.T) {
	t.Run("Test registering domains with success", func(t *testing.T) {
		if s := RegisterDomain(name, registerClientsSuccess); s != Owned {
			t.Logf("Expected %d result, got %d", Owned, s)
			t.Fail()
		}
	})
	t.Run("Test registering domains with failure", func(t *testing.T) {
		if s := RegisterDomain(name, registerClientsFailure); s != Unavailable {
			t.Logf("Expected %d result, got %d", Unavailable, s)
			t.Fail()
		}
	})
}

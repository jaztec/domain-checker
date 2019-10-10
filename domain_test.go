package checker

import (
	"testing"
)

const errorMessage = "some error message"

type expectedCheckDomainResult struct {
	registrar Registrar
	status    Status
}

var (
	domainRegistrars = []Registrar{
		availableRegistrar{},
		unavailableRegistrar{},
		ownedRegistrar{},
		processingRegistrar{},
		errorRegistrar{},
	}
	registerRegistrarsSuccess = []Registrar{
		errorRegistrar{},
		unavailableRegistrar{},
		ownedRegistrar{},
	}
	registerRegistrarsFailure = []Registrar{
		errorRegistrar{},
		unavailableRegistrar{},
	}
	name = "irrelevant"
)

func TestCheckDomain(t *testing.T) {
	t.Run("Test statuses", func(t *testing.T) {

		expectLen := len(domainRegistrars) - 1
		expectedResults := []expectedCheckDomainResult{
			{domainRegistrars[0], Available},
			{domainRegistrars[1], Unavailable},
			{domainRegistrars[2], Owned},
			{domainRegistrars[3], Processing},
		}

		statuses := CheckDomain(name, domainRegistrars)

		if gotLen := len(statuses); gotLen != expectLen {
			t.Logf("Expected %d result statuses but received %d", expectLen, gotLen)
			t.Fail()
		}

		// test the results in order of registrar ordering, it should comply
		for i, s := range expectedResults {
			status := statuses[i]
			if status.Status() != s.status {
				t.Logf("Expected %T from status but received %T", s.status, status.Status())
				t.Fail()
			}
			if status.Registrar() != s.registrar {
				t.Logf("Expected %T from registrar but received %T", s.registrar, status.Registrar())
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
		if s := RegisterDomain(name, registerRegistrarsSuccess); s.Status() != Owned {
			t.Logf("Expected %d result, got %d", Owned, s.Status())
			t.Fail()
		}
	})
	t.Run("Test registering domains with failure", func(t *testing.T) {
		if s := RegisterDomain(name, registerRegistrarsFailure); s.Status() != Unavailable {
			t.Logf("Expected %d result, got %d", Unavailable, s.Status())
			t.Fail()
		}
	})
}

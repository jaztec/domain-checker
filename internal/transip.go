package internal

import (
	"fmt"

	checker "github.com/jaztec/domain-checker"
	"github.com/transip/gotransip"
	transipDomain "github.com/transip/gotransip/domain"
)

type transip struct {
	client gotransip.Client
}

// CheckDomain will consult the TransIP services and return a modified internal Status on whether
// the domain is available for registration.
func (t *transip) CheckDomain(n string) (s checker.Status, err error) {
	s = checker.Unavailable
	ts, err := transipDomain.CheckAvailability(t.client, n)
	if err != nil {
		return s, fmt.Errorf("check domain availability returned an error: %w", checker.NewError(t, err))
	}

	switch ts {
	case transipDomain.StatusInYourAccount:
		s = checker.Owned
	case transipDomain.StatusInternalPush:
		s = checker.Owned
	case transipDomain.StatusFree:
		s = checker.Available
	}
	return
}

// RegisterDomain will try and register a certain domain name at the TransIP API.
func (t *transip) RegisterDomain(name string) (checker.Status, error) {
	err := transipDomain.Register(t.client, transipDomain.Domain{Name: name})
	if err != nil {
		return checker.Unavailable, err
	}
	return checker.Processing, nil
}

// NewTransIP returns a new client for site validations at TransIP
func NewTransIP(accountName, keyPath string) (checker.Registrar, error) {
	c, err := gotransip.NewSOAPClient(gotransip.ClientConfig{
		AccountName:    accountName,
		PrivateKeyPath: keyPath,
		Mode:           gotransip.APIModeReadWrite,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating TransIP client: %v", err)
	}
	t := &transip{&c}
	return t, nil
}

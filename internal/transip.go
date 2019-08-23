package internal

import (
	"github.com/jaztec/domain-checker/pkg/checker"
	"github.com/transip/gotransip"
	transipDomain "github.com/transip/gotransip/domain"
	"golang.org/x/xerrors"
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
		return s, xerrors.Errorf("check domain availability returned an error: %w", checker.NewError(t, err))
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

// NewTransIP returns a new client for site validations at TransIP
func NewTransIP(accountName, keyPath string) (checker.Client, error) {
	c, err := gotransip.NewSOAPClient(gotransip.ClientConfig{
		AccountName:    accountName,
		PrivateKeyPath: keyPath,
		Mode:           gotransip.APIModeReadWrite,
	})
	if err != nil {
		return nil, xerrors.Errorf("error creating TransIP client: %w", err)
	}
	t := &transip{&c}
	return t, nil
}
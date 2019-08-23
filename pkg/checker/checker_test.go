package checker

import "errors"

type availableClient struct{}
type unavailableClient struct{}
type ownedClient struct{}
type errorClient struct{}

func (availableClient) CheckDomain(string) (Status, error)      { return Available, nil }
func (availableClient) RegisterDomain(string) (Status, error)   { return Available, nil }
func (unavailableClient) CheckDomain(string) (Status, error)    { return Unavailable, nil }
func (unavailableClient) RegisterDomain(string) (Status, error) { return Unavailable, nil }
func (ownedClient) CheckDomain(string) (Status, error)          { return Owned, nil }
func (ownedClient) RegisterDomain(string) (Status, error)       { return Owned, nil }
func (errorClient) CheckDomain(string) (Status, error)          { return Unavailable, errors.New(errorMessage) }
func (errorClient) RegisterDomain(string) (Status, error) {
	return Unavailable, errors.New(errorMessage)
}

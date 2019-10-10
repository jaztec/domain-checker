package checker

import "errors"

type availableRegistrar struct{}
type unavailableRegistrar struct{}
type ownedRegistrar struct{}
type processingRegistrar struct{}
type errorRegistrar struct{}

func (availableRegistrar) CheckDomain(string) (Status, error)      { return Available, nil }
func (availableRegistrar) RegisterDomain(string) (Status, error)   { return Available, nil }
func (unavailableRegistrar) CheckDomain(string) (Status, error)    { return Unavailable, nil }
func (unavailableRegistrar) RegisterDomain(string) (Status, error) { return Unavailable, nil }
func (ownedRegistrar) CheckDomain(string) (Status, error)          { return Owned, nil }
func (ownedRegistrar) RegisterDomain(string) (Status, error)       { return Owned, nil }
func (processingRegistrar) CheckDomain(string) (Status, error)     { return Processing, nil }
func (processingRegistrar) RegisterDomain(string) (Status, error)  { return Processing, nil }
func (errorRegistrar) CheckDomain(string) (Status, error) {
	return Unavailable, errors.New(errorMessage)
}
func (errorRegistrar) RegisterDomain(string) (Status, error) {
	return Unavailable, errors.New(errorMessage)
}

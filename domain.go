package checker

import (
	"fmt"
	"log"
)

// Status wraps statuses this package will act upon
type Status uint8

var (
	// Unavailable status notes the domain name is not available to us in any feasible way
	Unavailable Status = 0x00
	// Owned status tells us the domain name is already in out possession
	Owned Status = 0x01
	// Available status is the mark that we can try to claim the domain name
	Available Status = 0x02
	// Processing is a special status that indicates we are already trying to get the domain
	Processing Status = 0x04
)

// ClientStatus tells the status for a domain for a specific domain
type RegistrarStatus struct {
	c      Registrar
	s      Status
	domain string
}

// Registrar reports the registrar to which this status applies
func (cs *RegistrarStatus) Registrar() Registrar {
	return cs.c
}

// Status reports the actual status for this domain with this client
func (cs *RegistrarStatus) Status() Status {
	return cs.s
}

// Domain reports the domain name requested
func (cs *RegistrarStatus) Domain() string {
	return cs.domain
}

// CheckDomain will walk though the provided domainRegistrars and check on all of them if a specific domain
// is available. The domainClients will be checked in order of appearance.
func CheckDomain(name string, clients []Registrar) []RegistrarStatus {
	results := make([]RegistrarStatus, 0, len(clients))

	for _, c := range clients {
		if s, err := c.CheckDomain(name); err == nil {
			results = append(results, RegistrarStatus{c, s, name})
		} else {
			log.Printf("%v", fmt.Errorf("received error from provider '%T' while checking domain '%s': %w", c, name, err))
		}
	}

	return results
}

// RegisterDomain will try to register a domain at a slice of given domainRegistrars. The first one to return a valid response
// will own the domain. Please sort the domainClients in order of preference.
func RegisterDomain(name string, clients []Registrar) (cs RegistrarStatus) {
	for _, c := range clients {
		if s, err := c.RegisterDomain(name); err == nil && (s == Owned || s == Processing) {
			cs = RegistrarStatus{
				c:      c,
				s:      s,
				domain: name,
			}
			return
		} else if err != nil {
			log.Printf("%v", fmt.Errorf("received error from provider '%T' while trying to register domain '%s': %w", c, name, err))
		}
	}
	return
}

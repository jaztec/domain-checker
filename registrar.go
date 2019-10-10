package checker

// Registrar interface defines some methods we want external services to present to us such as but not
// limited to domain availability checks and registration
type Registrar interface {
	// CheckDomain returns a status about the requested domain.
	CheckDomain(string) (Status, error)
	// Register domain will try and register the domain name
	RegisterDomain(string) (Status, error)
}

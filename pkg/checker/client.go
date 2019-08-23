package checker

// Client interface defines some methods we want external services to present to us such as but not
// limited to domain availability checks and registration
type Client interface {
	CheckDomain(string) (Status, error)
}

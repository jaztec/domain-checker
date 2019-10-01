package main

import (
	"log"
	"time"

	"github.com/jaztec/domain-checker/pkg/checker"
)

func loadClients() []checker.Client {
	return make([]checker.Client, 0)
}

func main() {
	var domains []string
	clients := loadClients()
	for {
		try := false
		for _, name := range domains {
			statuses := checker.CheckDomain(name, clients)
			for _, s := range statuses {
				if s.Status() == checker.Available {
					try = true
					break
				}
			}
			if try {
				if s := checker.RegisterDomain(name, clients); s.Status() == checker.Owned || s.Status() == checker.Processing {
					log.Printf("Registered '%s' at %T", name, s.Client())
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}

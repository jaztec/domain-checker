package main

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
	checker "github.com/jaztec/domain-checker"
)

type checking struct {
	redis      *redis.Client
	lock       sync.RWMutex
	domains    []string
	registrars []checker.Registrar
}

func (c *checking) runChecks() {
	for {
		c.lock.RLock()
		for _, name := range c.domains {
			statuses, _ := checker.CheckDomain(name, c.registrars)
			for _, s := range statuses {
				if s.Status() == checker.Available {
					if s, _ := checker.RegisterDomain(name, c.registrars); s.Status() == checker.Owned || s.Status() == checker.Processing {
						log.Printf("Registered '%s' at %T", name, s.Registrar())
					}
					break
				}
			}
		}
		c.lock.RUnlock()

		time.Sleep(60 * time.Second)
	}
}

func (c *checking) findDomain(name string) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for i, d := range c.domains {
		if name == d {
			return i
		}
	}
	return -1
}

func (c *checking) addDomain(name string) {
	if c.findDomain(name) == -1 {
		c.domains = append(c.domains, name)
	}
	c.persistRedis()
	log.Printf("Added domain \"%s\"", name)
}

func (c *checking) removeDomain(name string) {
	if i := c.findDomain(name); i != -1 {
		d := c.domains
		c.domains = d[:i+copy(d[i:], d[i+1:])]
	}
	c.persistRedis()
	log.Printf("Removed domain \"%s\"", name)
}

func (c *checking) listDomains() []string {
	return c.domains
}

func (c *checking) persistRedis() {
	if c.redis != nil {
		c.lock.Lock()
		defer c.lock.Unlock()
		pipe := c.redis.Pipeline()
		pipe.Del(RedisListKey)
		for _, d := range c.domains {
			pipe.LPush(RedisListKey, d)
		}
		pipe.Exec()
	}
}

func newChecking(domains []string, clients []checker.Registrar, r *redis.Client) *checking {
	return &checking{
		redis:   r,
		domains: domains,
	}
}

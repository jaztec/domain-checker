package main

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
	checker "github.com/jaztec/domain-checker"
)

type checking struct {
	redis   *redis.Client
	lock    sync.RWMutex
	domains []string
	clients []checker.Client
}

func (c *checking) runChecks() {
	for {
		c.lock.RLock()
		for _, name := range c.domains {
			statuses := checker.CheckDomain(name, c.clients)
			for _, s := range statuses {
				if s.Status() == checker.Available {
					if s := checker.RegisterDomain(name, c.clients); s.Status() == checker.Owned || s.Status() == checker.Processing {
						log.Printf("Registered '%s' at %T", name, s.Client())
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
}

func (c *checking) removeDomain(name string) {
	if i := c.findDomain(name); i != -1 {
		d := c.domains
		c.domains = d[:i+copy(d[i:], d[i+1:])]
	}
	c.persistRedis()
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

func newChecking(domains []string, clients []checker.Client, r *redis.Client) *checking {
	return &checking{
		redis:   r,
		domains: domains,
	}
}

package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	checker "github.com/jaztec/domain-checker"
	"github.com/jaztec/domain-checker/internal"
)

// RedisListKey defines the key within Redis that is used for
// caching the domain name list.
const RedisListKey = "checker_domain_list"

func startRedis(dsn, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       db,
	})

	p, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	if p != "PONG" {
		return nil, fmt.Errorf("%s is an unexpected PING result", p)
	}

	return client, nil
}

func loadClients() []checker.Registrar {
	c := make([]checker.Registrar, 0, 5)

	transIPName := os.Getenv("TRANSIP_ACCOUNT_NAME")
	transIPKey := os.Getenv("TRANSIP_KEY_FILE_PATH")
	if transIPName != "" && transIPKey != "" {
		internal.NewTransIP(transIPName, transIPKey)
		c = append(c)
	}

	return c
}

func main() {
	var domains []string

	// get Redis connection going to keep track of requested domains
	// over restarts. If no Redis connection can be established the
	// program will continue without persisten storage.
	var r *redis.Client
	dsn := os.Getenv("REDIS_DSN")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Printf("%v\n", fmt.Errorf("error while loading Redis db variable: %w", err))
	} else {
		r, err := startRedis(dsn, password, db)
		if err != nil {
			log.Printf("%v\n", (fmt.Errorf("error while conecting to Redis: %w", err)))
		} else {
			if n, err := r.LLen(RedisListKey).Result(); err != nil {
				var i int64
				for ; i < n; i++ {
					if s, err := r.RPop(RedisListKey).Result(); err == nil {
						domains = append(domains, s)
					}
				}
			}
			defer r.Close()
		}
	}

	// run the checking loops
	c := newChecking(domains, loadClients(), r)

	// get server running for communication with this instance
	port := os.Getenv("PORT")
	if port == "" {
		panic(errors.New("no valid port received to launch server"))
	}

	// create a configuration to setup a secure connection
	cert, err := tls.LoadX509KeyPair(os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY"))
	var tlsConfig *tls.Config
	if err != nil {
		if os.Getenv("TLS_ALLOW_INSECURE") != "true" {
			panic(err)
		}
		log.Println("WARNING: server is started without TLS, please make sure it is not connected to the internet of protected behind a proxy with SSL termination")
	} else {
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	token := os.Getenv("AUTH_TOKEN")
	s, err := newServer(port, token, c, tlsConfig)
	if err != nil {
		panic(fmt.Errorf("error while launching server: %w", err))
	}
	defer s.close()

	c.runChecks()
}

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

func isAuthenticated(c *client) bool {
	if !c.authenticated {
		c.write("unauthenticated")
	}
	return c.authenticated
}

type command struct {
	name   string
	params []string
}

type client struct {
	conn          net.Conn
	authenticated bool
	commands      chan command
}

func (c *client) close() {
	log.Printf("Closing connection to '%s'", c.conn.RemoteAddr())
	c.conn.Close()
	close(c.commands)
}

func (c *client) read() {
	for {
		d, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				return
			}
			log.Printf("an error occured with %s: %T", c.conn.RemoteAddr().String(), err)
			return
		}

		cmd := strings.Fields(string(d))
		if len(cmd) == 0 {
			continue
		}
		name := strings.ToUpper(cmd[0])
		params := make([]string, len(cmd)-1)
		for i := 1; i < len(cmd); i++ {
			params[i-1] = cmd[i]
		}
		c.commands <- command{
			name:   name,
			params: params,
		}
	}
}

func (c *client) write(s string) {
	if _, err := c.conn.Write([]byte(s + "\n")); err != nil {
		log.Printf("An error occured with connection '%s': %v", c.conn.RemoteAddr(), err)
	}
}

type server struct {
	listener net.Listener
	done     chan struct{}
	checking *checking
	token    string
}

func (s *server) close() {
	close(s.done)
	s.listener.Close()
}

func (s *server) loop() {
	for {
		select {
		case _ = <-s.done:
			return
		default:
			c, err := s.listener.Accept()
			if err != nil {
				continue
			}
			cl := &client{
				conn:          c,
				authenticated: false,
				commands:      make(chan command),
			}
			log.Printf("New connection from '%s'", c.RemoteAddr())
			go s.handle(cl)
		}
	}
}

func (s *server) handle(c *client) {
	defer c.close()
	go c.read()

	for {
		select {
		case _ = <-s.done:
			log.Println("'done' channel closed")
			return
		case cmd := <-c.commands:
			switch cmd.name {
			case "AUTH":
				if cmd.params[0] == s.token {
					c.authenticated = true
				} else {
					c.write("authentication failure")
				}
			case "ADD":
				if !isAuthenticated(c) {
					break
				}
				s.checking.addDomain(cmd.params[0])
				c.write(fmt.Sprintf("%s added", cmd.params[0]))
			case "REMOVE":
				if !isAuthenticated(c) {
					break
				}
				s.checking.removeDomain(cmd.params[0])
				c.write(fmt.Sprintf("%s removed", cmd.params[0]))
			case "LIST":
				if !isAuthenticated(c) {
					break
				}
				domains := s.checking.listDomains()
				res := ""
				for _, domain := range domains {
					res += domain + " "
				}
				c.write(res)
			case "EXIT", "QUIT", "CLOSE":
				return
			default:
				// ignore
			}
		}
	}
}

func newServer(port, token string, check *checking, tlsConf *tls.Config) (*server, error) {
	var l net.Listener
	var err error
	if tlsConf == nil {
		l, err = net.Listen("tcp", ":"+port)
	} else {
		l, err = tls.Listen("tcp", ":"+port, tlsConf)
	}
	if err != nil {
		return nil, fmt.Errorf("an error occured: %w", err)
	}
	s := &server{
		listener: l,
		done:     make(chan struct{}, 1),
		checking: check,
		token:    token,
	}
	go s.loop()

	return s, nil
}

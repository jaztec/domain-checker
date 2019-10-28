package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	listener net.Listener
	done     chan struct{}
	checking *checking
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
			go s.handle(c)
		}
	}
}

func (s *server) handle(c net.Conn) {
	defer c.Close()

	// TODO This code will be blocking while reading the connection stream, it should be moved to another goroutine
	for {
		select {
		case _ = <-s.done:
			return
		default:
			d, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				log.Printf("an error occured with %s: %v", c.RemoteAddr().String(), err)
				return
			}

			cmd := strings.Fields(string(d))
			if len(cmd) == 0 {
				continue
			}
			log.Printf("Process command '%v'", cmd)
			switch strings.ToUpper(cmd[0]) {
			case "ADD":
				if len(cmd) < 2 {
					break
				}
				s.checking.addDomain(cmd[1])
			case "REMOVE":
				if len(cmd) < 2 {
					break
				}
				s.checking.removeDomain(cmd[1])
			case "LIST":
				domains := s.checking.listDomains()
				for _, domain := range domains {
					_, _ = c.Write([]byte(domain + " "))
				}
				c.Write([]byte("\n"))
			case "EXIT":
			case "QUIT":
			case "CLOSE":
				log.Printf("Closing connection from \"%s\"", c.RemoteAddr().String())
				_, _ = c.Write([]byte("Closing connection\n"))
				c.Close()
				return
			default:
				// ignore
			}
		}
	}
}

func newServer(port string, check *checking, tlsConf *tls.Config) (*server, error) {
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
	}
	go s.loop()

	return s, nil
}

package main

import (
	"bufio"
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

	for {
		select {
		case _ = <-s.done:
			return
		default:
			d, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				log.Printf("an error occured: %v", err)
				return
			}

			cmd := strings.Fields(string(d))
			switch cmd[0] {
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
			case "STOP":
				s.close()
				return
			default:
				// ignore
			}
		}
	}
}

func newServer(port string, check *checking) (*server, error) {
	l, err := net.Listen("tcp", ":"+port)
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

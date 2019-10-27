package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/urfave/cli"
)

var done = make(chan struct{}, 1)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	f := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "localhost",
			Usage: "Set the hostname of the domain-checker server",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 8081,
			Usage: "Set the port of the domain-checker server",
		},
		cli.BoolFlag{
			Name:  "allow-unsafe",
			Usage: "Use this flag to allow unsafe TLS connections",
		},
		cli.BoolFlag{
			Name:  "force-unsafe",
			Usage: "Use this flag to force a regular connection",
		},
	}

	getConn := func(c *cli.Context) net.Conn {
		conn, err := createConnection(c.String("host"), c.Int("port"), c.Bool("force-unsafe"), c.Bool("allow-unsafe"))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return conn
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all domains the server is currently watching",
			Flags:   f,
			Action: func(c *cli.Context) error {
				var wg sync.WaitGroup
				conn := getConn(c)
				defer closeConnection(conn)
				wg.Add(1)
				go readFromConnection(conn, &wg)

				_, err := conn.Write([]byte("LIST\n"))
				if err != nil {
					return err
				}

				wg.Wait()

				return nil
			},
		},
		cli.Command{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Use 'add [domain]' to add a domain to the checker list",
			Flags:   f,
			Action: func(c *cli.Context) error {
				if len(c.Args()) == 0 {
					return errors.New("no domain name provided")
				}
				if len(c.Args()) > 1 {
					return fmt.Errorf("too many parameters received: %v", c.Args())
				}
				domain := c.Args()[0]
				if len(domain) > 255 {
					return fmt.Errorf("domain name contains too many characters: %s", domain)
				}
				conn := getConn(c)
				defer closeConnection(conn)

				_, err := conn.Write([]byte("ADD " + domain + "\n"))
				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createConnection(host string, port int, forceUnsafe, allowUnsafe bool) (net.Conn, error) {
	var conn net.Conn
	var err error
	conf := tls.Config{}
	if allowUnsafe {
		conf.InsecureSkipVerify = true
	}
	if forceUnsafe {
		conn, err = net.Dial("tcp", host+":"+strconv.Itoa(port))
	} else {
		conn, err = tls.Dial("tcp", host+":"+strconv.Itoa(port), &conf)
	}
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func readFromConnection(conn net.Conn, wg *sync.WaitGroup) {
	ch := make(chan string)
	go func() {
		d, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("an error occured with %s: %v", conn.RemoteAddr().String(), err)
			return
		}

		ch <- fmt.Sprintf("%s", string(d))
	}()
	for {
		select {
		case _ = <-done:
			close(ch)
			return
		case s := <-ch:
			fmt.Println(s)
			wg.Done()
		}
	}
}

func closeConnection(conn net.Conn) {
	close(done)
	_, err := conn.Write([]byte("CLOSE\n"))
	if err != nil {
		log.Fatalf("Fatal exception occured: %v", err)
	}
}

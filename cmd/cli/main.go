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
	"time"

	"github.com/urfave/cli"
)

var done = make(chan struct{}, 1)

func doCommand(c net.Conn, cmd string) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	f := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "localhost",
			Usage: "Set the hostname of the domain-checker server",
		},
		cli.StringFlag{
			Name:  "token",
			Value: "",
			Usage: "Set the token to use to connect to the domain-checker server",
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

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all domains the server is currently watching with 'list'",
			Flags:   f,
			Action: func(c *cli.Context) error {
				conn, _ := getConn(c)
				defer closeConnection(conn)
				if err := doCommand(conn, "LIST\n"); err != nil {
					return err
				}
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
				conn, _ := getConn(c)
				defer closeConnection(conn)

				if err := doCommand(conn, "ADD "+domain+"\n"); err != nil {
					return err
				}
				return nil
			},
		},
		cli.Command{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "Use 'remove [domain]' to add a domain to the checker list",
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
				conn, _ := getConn(c)
				defer closeConnection(conn)

				if err := doCommand(conn, "REMOVE "+domain+"\n"); err != nil {
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

func getConn(c *cli.Context) (net.Conn, error) {
	conn, err := createConnection(c.String("host"), c.Int("port"), c.Bool("force-unsafe"), c.Bool("allow-unsafe"))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	go readFromConnection(conn)
	auth := "AUTH " + c.String("token") + "\n"
	if _, err := conn.Write([]byte(auth)); err != nil {
		return nil, err
	}
	return conn, nil
}

func readFromConnection(conn net.Conn) {
	ch := make(chan string)
	defer close(ch)
	go func() {
		d, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("an error occured with %s: %v", conn.RemoteAddr(), err)
			return
		}

		ch <- fmt.Sprintf("%s", string(d))
	}()
	for {
		select {
		case _ = <-done:
			return
		case s := <-ch:
			fmt.Println(s)
			return
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

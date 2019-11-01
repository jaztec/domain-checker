package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type config struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Token       string `yaml:"token"`
	AllowUnsafe bool   `yaml:"allowUnsafe"`
	ForceUnsafe bool   `yaml:"forceUnsafe"`
}

var (
	done       = make(chan struct{}, 1)
	configFile string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	configFile = usr.HomeDir + "/.checker-config.yaml"
}

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

	cfg := config{
		Host: "localhost",
		Port: 8081,
	}

	// Load the config file if present
	if f, err := ioutil.ReadFile(configFile); err == nil {
		err := yaml.Unmarshal(f, &cfg)
		if err != nil {
			panic(err)
		}
	}

	// Make sure env variables are correctly set on bool flags if the config requires them to be set
	if cfg.AllowUnsafe {
		os.Setenv("CHECKER_ALLOW_UNSAFE", "true")
	} else {
		os.Setenv("CHECKER_ALLOW_UNSAFE", "false")
	}
	if cfg.ForceUnsafe {
		os.Setenv("CHECKER_FORCE_UNSAFE", "true")
	} else {
		os.Setenv("CHECKER_FORCE_UNSAFE", "false")
	}

	f := []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: cfg.Host,
			Usage: "Set the hostname of the domain-checker server",
		},
		cli.StringFlag{
			Name:  "token",
			Value: cfg.Token,
			Usage: "Set the token to use to connect to the domain-checker server",
		},
		cli.IntFlag{
			Name:  "port",
			Value: cfg.Port,
			Usage: "Set the port of the domain-checker server",
		},
		cli.BoolFlag{
			Name:   "allow-unsafe",
			Usage:  "Use this flag to allow unsafe TLS connections",
			EnvVar: "CHECKER_ALLOW_UNSAFE",
		},
		cli.BoolFlag{
			Name:   "force-unsafe",
			Usage:  "Use this flag to force a regular connection",
			EnvVar: "CHECKER_FORCE_UNSAFE",
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
		cli.Command{
			Name:  "set",
			Usage: "Use 'set [name] [value]' to persist variables to the cli tool config",
			Action: func(c *cli.Context) error {
				return updateConfig(&cfg, c.Args())
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

func updateConfig(cfg *config, args cli.Args) error {
	// validate exact number of arguments
	if len(args) != 2 {
		return errors.New("Invalid arguments provided")
	}

	n := strings.Title(args[0])
	v := args[1]

	// check if the field argument exists in our configuration model
	t := reflect.ValueOf(cfg).Elem()
	_, found := t.Type().FieldByName(n)
	if !found {
		return fmt.Errorf("'%s' is not a valid configuration field", n)
	}

	// transform the type of the value to be usable for our configuration model
	var value interface{}
	if v == "true" {
		value = true
	} else if v == "false" {
		value = false
	} else if tmp, err := strconv.Atoi(v); err == nil {
		value = tmp
	} else {
		// it's not a bool, it's not an int, it's your problem now
		value = v
	}
	vT := reflect.TypeOf(value)
	if vT.Kind() != t.FieldByName(n).Type().Kind() {
		return fmt.Errorf("'%s' is not a valid configuration value for '%s'. Expected '%s'", vT.Kind(), n, t.FieldByName(n).Type().Kind())
	}

	// the config option exists and the value is of the correct type
	fT := t.FieldByName(n)
	if !fT.CanSet() {
		return fmt.Errorf("field '%s' is not settable", n)
	}
	fT.Set(reflect.ValueOf(value))

	// write the configuration to disk or return the failure
	return writeConfig(cfg)
}

func writeConfig(cfg *config) (err error) {
	var b []byte
	if b, err = yaml.Marshal(*cfg); err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, b, 0644)
}

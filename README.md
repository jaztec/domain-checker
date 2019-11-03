[![Build Status](https://travis-ci.com/jaztec/domain-checker.svg?branch=master)](https://travis-ci.com/jaztec/domain-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaztec/domain-checker)](https://goreportcard.com/report/github.com/jaztec/domain-checker)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/jaztec/domain-checker/blob/master/LICENSE)
[![GoDoc Domain Checker](https://godoc.org/github.com/jaztec/domain-checker?status.svg)](https://godoc.org/github.com/jaztec/domain-checker)
[![](https://images.microbadger.com/badges/image/jaztec/domain-checker.svg)](https://microbadger.com/images/jaztec/domain-checker)
[![](https://images.microbadger.com/badges/version/jaztec/domain-checker.svg)](https://microbadger.com/images/jaztec/domain-checker)

# Domain Checker

The purpose of this repository is to setup a list of domains you would like to
own. The program will check the domain with the providers available and to which
you have an account to register it on availability and will register it when it 
becomes available.

## Howto

### How to use the server
Use the `.env.dist` and the `docker-compose.yml.example.dist` to set up a running server.
You can add, remove and list domains that are being checked with the internal CLI command, 
the server container gets packed with the CLI program as well. You can use them by calling
the program from inside the container with `$ docker-compose exec checker cli list`.
You should also set the `AUTH_TOKEN` environment variable. New client connections need to 
authenticate with this token before any actions can be performed.

#### TLS
When you provide the server with a TLS cert and key the server will start a TLS server to
handle incoming connections. When no TLS vars are set or the vars are invalid the server
will try to fallback on a regular TCP server if the environment variable TLS_ALLOW_INSECURE
allows it.

#### Redis
You can add a Redis connection string to the server environment variables, this will make
sure the list of domains gets persisted into Redis so it can resume where it was after a
restart. If no Redis parameters are provided the server application will just keep the list
in memory.

### How to use the CLI program
The CLI program is packed with the server program into one Docker container. However it is 
also possible to use the CLI program standalone on a different computer. You can download this
repository and `make build` all programs inside, the CLI program will be located in 
`./bin/cli`.

#### Arguments
The client needs multile input arguments to setup a connection with the Checker service, but it
is also possible to set these values with the program so they are used automatically at each
command. The arguments available are:

Name | Workings | Required
--- | --- | ---
`host` | The hostname to connect to | True
`port` | The port on the host to connect to | True
`token` | The token to authenticate to the server with | True
`allowUnsafe` | Allow TLS connections with an invalid root certificate (dev mode) | False
`forceUnsafe` | Use a regular TCP connection, this should be used internally only and only works if the server isn't using TLS | False

All these variables can be persisted to the client app as well, use `$ cli set host example.org`, or `$ cli set forceUnsafe true` to set the variable with the client. The next
time the client program is called you no longer need to provide the flag to the cli command 
itself.

#### Commands
The application accepts 3 commands, `add`, `remove` and `list`. You can use them as follows
`$ cli [arguments] add host.com`. Or `$ cli [arguments] list`.

## Roadmap
- It would be nice if the server and CLI command do some domain name validation before adding/removing them.

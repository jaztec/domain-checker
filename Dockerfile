FROM golang:1.13-alpine as builder
LABEL maintainer="Jasper van Herpt <jasper.v.herpt@gmail.com>"

# Create user for the app
RUN adduser -D -g '' app-user

ENV GO111MODULE=on

WORKDIR /opt/local
COPY . .

RUN apk update && apk add --no-cache git make ca-certificates && update-ca-certificates

RUN go mod vendor \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./bin/checker ./cmd/checker \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./bin/cli ./cmd/cli

FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /opt/local/bin/checker /usr/bin/checker
COPY --from=builder /opt/local/bin/cli /usr/bin/checker-cli

USER app-user

CMD ["/usr/bin/checker"]

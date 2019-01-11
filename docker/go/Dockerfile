FROM golang:1.11-alpine3.8

RUN apk add --no-cache curl gcc git libc-dev
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
    | sh -s -- -b $GOPATH/bin v1.12.5

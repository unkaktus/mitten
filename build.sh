#!/bin/bash
set -e

PROGRAM=mitten
VERSION=$(git describe --exact-match --tags)
echo $VERSION
LDFLAGS="-X main.version=$VERSION"

env GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o $PROGRAM-Linux-x86_64
env GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o $PROGRAM-Linux-aarch64
env GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o $PROGRAM-Darwin-x86_64
env GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o $PROGRAM-Darwin-arm64
env GOOS=openbsd GOARCH=amd64 go build -ldflags "$LDFLAGS" -o $PROGRAM-OpenBSD-amd64
env GOOS=openbsd GOARCH=arm64 go build -ldflags "$LDFLAGS" -o $PROGRAM-OpenBSD-arm64

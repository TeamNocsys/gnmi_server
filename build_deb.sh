#!/bin/sh

set -e

#GOOS=linux go build -gcflags "all=-N -l" -tags release -o build/deb/usr/local/bin cmd/gnmi/gnmi.go
dpkg -b build/deb gnmi.deb
dpkg -c gnmi.deb

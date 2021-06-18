#!/bin/bash

set -e

GOOS=linux go build -gcflags "all=-N -l" -tags release -o build/deb/usr/local/bin cmd/gnmi/gnmi.go

chmod +x build/deb/usr/local/bin/gnmi
dpkg -b build/deb gnmi.deb
dpkg -c gnmi.deb

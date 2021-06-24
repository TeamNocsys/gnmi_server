#!/bin/sh

set -e

GIT_COMMIT_SHORT=$(git rev-parse --short HEAD)

sed -i "s/Version:/Version: 1.$GIT_COMMIT_SHORT/g" build/deb/DEBIAN/control
#GOOS=linux go build -gcflags "all=-N -l" -tags release -o build/deb/usr/local/bin cmd/gnmi/gnmi.go
dpkg -b build/deb gnmi.deb
dpkg -c gnmi.deb

#!/bin/bash

set -e

# put git commit hash info into debian control file
CPATH=build/deb/DEBIAN/control
git checkout -- $CPATH
GIT_COMMIT=$(git describe --dirty --always)
sed -i 's/GIT_COMMIT/'$GIT_COMMIT'/g' $CPATH

GOOS=linux go build -gcflags "all=-N -l" -tags release -o build/deb/usr/local/bin cmd/gnmi/gnmi.go

chmod +x build/deb/usr/local/bin/gnmi
dpkg -b build/deb gnmi.deb
dpkg -c gnmi.deb

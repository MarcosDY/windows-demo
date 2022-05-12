#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

SPIRE_RELEASE=1.3.0
SPIRE_URL=https://github.com/spiffe/spire/releases/download/v${SPIRE_RELEASE}/spire-${SPIRE_RELEASE}-windows-x86_64.zip

echo "Building containers and binaries"
./build.sh

if [ ! -d "spire" ]; then
   echo "Downloading SPIRE release"
   curl -o spire.zip -sSfL $SPIRE_URL
   echo "Installing SPIRE"
   unzip spire.zip
   rm spire.zip
   mv spire-$SPIRE_RELEASE/ spire/
fi

echo "Updating SPIRE config"
cp spire-conf/agent/* spire/conf/agent
cp spire-conf/server/* spire/conf/server

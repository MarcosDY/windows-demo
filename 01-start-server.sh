#!/bin/bash

set -e

(cd spire; ./bin/spire-server.exe run -config conf/server/server.conf)

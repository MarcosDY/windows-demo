#!/bin/bash

set -e

# TODO: how can we start as admin from here?
(cd spire; ./bin/spire-server.exe run -config conf/server/server.conf)

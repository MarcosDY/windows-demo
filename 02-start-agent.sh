#!/bin/bash

set -e

# TODO: how can we start as admin from here?
(cd spire; ./bin/spire-agent.exe run -config conf/agent/agent.conf)

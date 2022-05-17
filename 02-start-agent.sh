#!/bin/bash

set -e

(cd spire; ./bin/spire-agent.exe run -config conf/agent/agent.conf)

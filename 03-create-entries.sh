#!/bin/bash

set -e

fingerprint() {
	# calculate the sha1 digest of the der bytes of the certificate using the
	# "coreutils" output format (`-r`) to provide uniform output from
	# `openssl sha1` on macos and linux.
	cat $1 | openssl x509 -outform der | openssl sha1 -r | awk '{print $1}'
}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
AGENT_FINGERPRINT=$(fingerprint "${DIR}"/spire/conf/agent/agent.crt.pem)
USER_NAME="$USERDOMAIN\\$USERNAME"

./spire/bin/spire-server.exe entry create \
     -parentID "spiffe://example.org/spire/agent/x509pop/${AGENT_FINGERPRINT}" \
     -spiffeID "spiffe://example.org/webapp" \
     -selector "docker:label:com.docker.compose.service:webapp" \
     -selector "docker:image_id:webapp" \
     -ttl 60

./spire/bin/spire-server.exe entry create \
     -parentID "spiffe://example.org/spire/agent/x509pop/${AGENT_FINGERPRINT}" \
     -spiffeID "spiffe://example.org/customers-api" \
     -selector "docker:label:com.docker.compose.service:customers-api" \
     -selector "docker:image_id:customer-api" \
     -ttl 60

./spire/bin/spire-server.exe entry create \
     -parentID "spiffe://example.org/spire/agent/x509pop/${AGENT_FINGERPRINT}" \
     -spiffeID "spiffe://example.org/products-api" \
     -selector "windows:user_name:$USER_NAME" \
     -ttl 60

#!/bin/bash

set -eu

fingerprint() {
	# calculate the sha1 digest of the der bytes of the certificate using the
	# "coreutils" output format (`-r`) to provide uniform output from
	# `openssl sha1` on macos and linux.
	cat $1 | openssl x509 -outform der | openssl sha1 -r | awk '{print $1}'
}

AGENT_FINGERPRINT=$(fingerprint ./spire/conf/agent/agent.crt.pem)

ENTRY_ID=$(./spire/bin/spire-server entry show -spiffeID spiffe://example.org/products-api | grep 'Entry ID         :' | awk '{print $4}')
USER_NAME="$USERDOMAIN\\$USERNAME"

./spire/bin/spire-server entry update \
     -entryID $ENTRY_ID \
     -parentID "spiffe://example.org/spire/agent/x509pop/${AGENT_FINGERPRINT}" \
     -spiffeID "spiffe://example.org/products-api-u" \
     -selector "windows:user_name:$USER_NAME" \
     -ttl 60

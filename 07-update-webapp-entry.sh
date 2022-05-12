#!/bin/bash

set -eu

fingerprint() {
	# calculate the sha1 digest of the der bytes of the certificate using the
	# "coreutils" output format (`-r`) to provide uniform output from
	# `openssl sha1` on macos and linux.
	cat $1 | openssl x509 -outform der | openssl sha1 -r | awk '{print $1}'
}

AGENT_FINGERPRINT=$(fingerprint ./spire/conf/agent/agent.crt.pem)

ENTRY_ID=$(./spire/bin/spire-server entry show -spiffeID spiffe://example.org/webapp | grep 'Entry ID         :' | awk '{print $4}')

./spire/bin/spire-server entry update \
     -entryID $ENTRY_ID \
     -parentID "spiffe://example.org/spire/agent/x509pop/${AGENT_FINGERPRINT}" \
     -spiffeID "spiffe://example.org/webapp" \
     -selector "docker:label:com.docker.compose.service:webapp" \
     -selector "docker:image_id:webapp" \
     -ttl 60

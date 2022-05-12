#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Building customer-api"
(cd "${DIR}"/src/customer-api && CGO_ENABLED=0 GOOS=windows go build -v -o $DIR/docker/customer-api/customerAPI.exe)
echo "Building webapp"
(cd "${DIR}"/src/webapp && CGO_ENABLED=0 GOOS=windows go build -v -o $DIR/docker/webapp/webapp.exe)
echo "Building product-api"
(cd "${DIR}"/src/product-api && CGO_ENABLED=0 GOOS=windows go build -v -o $DIR/productAPI.exe)

docker-compose -f "${DIR}"/docker-compose.yaml build

#!/usr/bin/env bash

# This script is used to build the batch service.

set -ex -o pipefail

# go the root of the project 
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
pushd "${SCRIPT_DIR}" > /dev/null || exit 1


SHA="$(git rev-parse --short HEAD)"
DATE="$(TZ=UTC date +"%Y-%m-%dT%TZ")"
LDFLAGS="-X main.gitSha=${SHA} -X main.buildTime=${DATE}"

go build \
    -a \
    -v \
    -trimpath='true' \
    -buildmode='exe' \
    -buildvcs='true' \
    -compiler='gc' \
    -mod='vendor' \
    -ldflags "${LDFLAGS}" \
    -o bin/batch

echo
echo "binaries built: bin/batch"
echo

# OR build the container image
# docker build -t ${REGION}.pkg.dev/${PROJECT_ID}/${SERVICE_NAME}/${SERVICE_NAME}:${SHA} .

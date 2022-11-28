#!/usr/bin/env bash
set -e

GOIMPORTS="$(go env GOPATH)/bin/goimports"

if [[ ! -x "${GOIMPORTS}" ]]; then
    echo "Installing goimports ..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

"${GOIMPORTS}" -w -local github.com/tomcz/gotools .

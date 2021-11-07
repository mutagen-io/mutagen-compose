#!/bin/bash

# Exit immediately on failure.
set -e

# Run tests and generate a coverage profile.
go test -v -coverpkg=./pkg/... -coverprofile=coverage.txt ./pkg/...

# Run tests with the race detector enabled. We use a slim end-to-end test since
# the race detector significantly increases the execution time.
go test -race ./pkg/...

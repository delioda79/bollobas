#!/bin/sh
set -e

# lint
echo "Checking lint"
golint -set_exit_status=1 `go list ./...`
echo "Lint success!"

# test
echo "Running tests"
go test `go list ./... | grep -v 'docs'` -mod=vendor -race -cover -tags=integration -coverprofile=coverage.txt -covermode=atomic
echo "Test execution completed!"
#!/usr/bin/env sh
foldersToLint=$(go list ./... | grep -v vendor)
golint -set_exit_status=1 ${foldersToLint}

#!/bin/sh

export LAMP_CONFIG_PATH=$PWD/contrib

go test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout 30s
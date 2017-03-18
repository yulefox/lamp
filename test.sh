#!/bin/sh
set -e

export LAMP_CONFIG_PATH=$PWD/contrib

go test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout 30s github.com/yulefox/lamp/core
go test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout 30s github.com/yulefox/lamp/apps

#!/bin/sh

export LAMP_CONFIG_PATH=$PWD/contrib

go test -v -timeout 300s
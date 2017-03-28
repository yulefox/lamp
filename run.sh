#!/bin/sh
set -e

file=$1
dir=`dirname $file`
export LAMP_CONFIG_PATH=$PWD/$dir

go run $file
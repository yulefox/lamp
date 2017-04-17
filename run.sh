#!/bin/sh
set -e

in=$1
dir=`dirname $in`
export LAMP_CONFIG_PATH=$PWD/$dir

if [ -f $in ] ; then
    go run $in
else
    go test $in
fi
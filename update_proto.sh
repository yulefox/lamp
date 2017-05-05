#! /bin/bash -e

dir=kits/cli
protoc -I ${dir} ${dir}/*.proto --go_out=plugins=grpc:${dir}/proto

echo >&2 "update protos OK"
#!/bin/sh

curl -l -H "Content-type: application/json" -X POST -d '{"host":"yulefox","data":"hello, world"}' http://localhost:4151/pub?topic=cmd
curl -l -H "Content-type: application/json" -X POST -d '{"host":"yulefox","data":"hello, world"}' http://localhost:4151/pub?topic=api

#!/bin/sh

curl -l -H "Content-type: application/json" -X POST -d '{"host":"local_mac","data":"hello, world"}' http://10.8.2.74:4151/pub?topic=gm

#! /bin/bash -e

root=github.com/yulefox/lamp
dir=kits
pkg=something
interfaces=MyInterface

mockgen ${root}/${dir}/${pkg} ${interfaces} \
  > ${dir}/${pkg}/mock_${pkg}/mock_${pkg}.go

echo >&2 "update mocks OK"
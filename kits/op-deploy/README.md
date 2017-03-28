# 运维-部署

## Install

```sh
go get github.com/yulefox/lamp/kits/op-deploy
```

## Usage

```sh
op-deploy 
```

## 流程

- 通过 `ssh` 执行初始化脚本

```sh
ssh -p<port> <user>@<hostname> 'wget http://lamp.yulefox.com/bootstrap.sh && bash bootstrap.sh'
```


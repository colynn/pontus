# Pontus

海神-蓬托斯

## Prerequisites

* `go` 1.18+
* `make`

## How to develop

### 1. 安装依赖

> 项目使用 go mod的方式来管理依赖包

```sh
    go mod vendor
```

### 2. 确认配置条件

* 默认会使用 ```config/local.yaml```， 请基于`dev.yaml`自行拷贝，依赖环境变量来获取不同配置文件

* __注意__: 本地测试时，根据需求修改`db`等相关配置

### 3. 启动-本地

```sh
    export ENV=local
    make run
```

### 4. 启动-生产环境

```sh
    export ENV=prod
    make run-prod
```

## Unit Test

### 单个包运行单元测试

```sh
go test -v [pkg-name] // eg， go test -v newops/pkg/mail
```

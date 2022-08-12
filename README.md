# gin-connect-demo

[![Code Generator](https://github.com/go-training/gin-connect-demo/actions/workflows/codegen.yaml/badge.svg)](https://github.com/go-training/gin-connect-demo/actions/workflows/codegen.yaml)

An example service built with [Connect](https://connect.build) and Gin.

## Install tools

```sh
make install
```

## Build Server and Client

```sh
make build
```

## Start server

### [Chi](https://github.com/go-chi/chi)

```sh
$ ./bin/chi-server
2022/08/12 10:05:00 Got connection: HTTP/2.0
2022/08/12 10:05:00 Request headers:  map[Accept-Encoding:[gzip] Content-Type:[application/proto] Gitea-Header:[hello from connect] User-Agent:[connect-go/0.2.0 (go1.19)]]
2022/08/12 10:05:00 "POST http://localhost:8080/proto.v1.GiteaService/Gitea HTTP/2.0" from [::1]:58227 - 200 40B in 1.36225ms
```

### [Gin](https://github.com/gin-gonic/gin)

```sh
$ ./bin/gin-server
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /proto.v1.GiteaService/:name --> main.giteaHandler.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080
```

## Start client

```sh
$ ./bin/client
2022/08/12 10:05:00 Message: Hello, foobar!
2022/08/12 10:05:00 Gitea-Version: v1
```

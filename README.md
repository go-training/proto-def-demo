# gin-connect-demo

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

```sh
$ ./bin/server
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
./bin/client
```

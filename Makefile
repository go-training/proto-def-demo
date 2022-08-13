GO ?= go

.PHONY: build
build: generator server client

.PHONY: server
server: gin chi

.PHONY: chi
chi:
	$(GO) build -o bin/$@-server cmd/server/$@/*.go

.PHONY: gin
gin:
	$(GO) build -o bin/$@-server cmd/server/$@/*.go

.PHONY: client
client:
	$(GO) build -o bin/$@ cmd/$@/main.go

.PHONY: install
install:
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest
	$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	$(GO) get -u -t ./... && go mod tidy -v

.PHONY: generator
generator:
	buf lint && buf generate

clean:
	rm -rf gen

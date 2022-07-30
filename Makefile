GO ?= go

install:
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest
	$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

build: server client

server:
	$(GO) build -o bin/$@ cmd/$@/main.go

client:
	$(GO) build -o bin/$@ cmd/$@/main.go

generator:
	buf lint && buf generate

clean:
	rm -rf gen

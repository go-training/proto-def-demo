
# for buf
BUF_BINARY_NAME=buf
BUF_VERSION=v1.7.0
OS := $(shell uname -s)-$(shell uname -m)
BIN=$(shell go env GOPATH)/bin

install:
	curl -sSL \
	  "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/${BUF_BINARY_NAME}-${OS}" \
	  -o ${BIN}/${BUF_BINARY_NAME} && \
	chmod +x ${BIN}/${BUF_BINARY_NAME}

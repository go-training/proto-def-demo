GO ?= go
BUF_VERSION=v1.7.0
GRPCURL_VERSION=v1.8.7
PROTOC_GEN_GO=v1.28
PROTOC_GEN_GO_GRPC=v1.2
PROTOC_GEN_CONNECT_GO=v0.4.0
PROTOC_GEN_OPENAPIV2=v2.11.3
PROTO_GO_TARGET_REPO ?= deploy/proto-go
PROTO_PYTHON_TARGET_REPO ?= deploy/proto-python
PROTO_RUBY_TARGET_REPO ?= deploy/proto-ruby
PROTO_OPENAPIV2_TARGET_REPO ?= deploy/proto-openapiv2

.PHONY: build
build: generator

.PHONY: install
install:
	$(GO) install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@$(GRPCURL_VERSION)
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO)
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC)
	$(GO) install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@$(PROTOC_GEN_CONNECT_GO)
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@$(PROTOC_GEN_OPENAPIV2)

.PHONY: generator
generator: buf-lint buf-gen-go buf-gen-python buf-gen-openapiv2 buf-gen-ruby

.PHONY: buf-lint
buf-lint:
	buf lint
	buf format --diff --exit-code

.PHONY: buf-format
buf-format:
	buf format --diff -w

buf-gen-go: clean_go
	buf generate --template buf.gen-go.yaml

buf-gen-ruby: clean_ruby
	buf generate --template buf.gen-ruby.yaml

buf-gen-python: clean_python
	buf generate --template buf.gen-python.yaml

buf-gen-openapiv2: clean_openapiv2
	buf generate --template buf.gen-openapiv2.yaml

push-to-go-repo:
	cp -r gen/go/* $(PROTO_GO_TARGET_REPO)/
	cd $(PROTO_GO_TARGET_REPO) && $(GO) mod init github.com/go-training/proto-go-demo || true
	cd $(PROTO_GO_TARGET_REPO) && $(GO) mod tidy
	git config --global user.email "appleboy.tw@gmail.com"
	git config --global user.name "Bo-Yi Wu"
	(cd $(PROTO_GO_TARGET_REPO) && git add --all && git commit -m "[auto-commit] Generate codes" && git push -f -u origin main) || echo "not pushed"

push-to-python-repo:
	cp -r gen/python/* $(PROTO_PYTHON_TARGET_REPO)/
	git config --global user.email "appleboy.tw@gmail.com"
	git config --global user.name "Bo-Yi Wu"
	(cd $(PROTO_PYTHON_TARGET_REPO) && git add --all && git commit -m "[auto-commit] Generate codes" && git push -f -u origin main) || echo "not pushed"

push-to-ruby-repo:
	cp -r gen/ruby/* $(PROTO_RUBY_TARGET_REPO)/
	git config --global user.email "appleboy.tw@gmail.com"
	git config --global user.name "Bo-Yi Wu"
	(cd $(PROTO_RUBY_TARGET_REPO) && git add --all && git commit -m "[auto-commit] Generate codes" && git push -f -u origin main) || echo "not pushed"

push-to-openapiv2-repo:
	cp -r gen/openapiv2/* $(PROTO_OPENAPIV2_TARGET_REPO)/
	git config --global user.email "appleboy.tw@gmail.com"
	git config --global user.name "Bo-Yi Wu"
	(cd $(PROTO_OPENAPIV2_TARGET_REPO) && git add --all && git commit -m "[auto-commit] Generate codes" && git push -f -u origin main) || echo "not pushed"

clean_go:
	rm -rf gen/go

clean_python:
	rm -rf gen/python

clean_openapiv2:
	rm -rf gen/openapiv2

clean_ruby:
	rm -rf gen/ruby

clean:
	rm -rf gen

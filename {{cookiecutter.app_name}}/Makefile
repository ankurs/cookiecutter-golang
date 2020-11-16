.PHONY: build build-alpine clean test help default proto install-deps


BIN_NAME={{cookiecutter.app_name}}

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

default: test

help:
	@echo 'Management commands for {{cookiecutter.app_name}}:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs dep ensure, mostly used for ci.'
	@echo '    make build-alpine    Compile optimized for alpine linux.'
	@echo '    make package         Build final docker image with just the go binary inside'
	@echo '    make tag             Tag image created by package with latest, git commit and version'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make push            Push tagged images to registry'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X {{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version.BuildDate=${BUILD_DATE}" -o bin/${BIN_NAME}

deps:
	go mod tidy

install:
	go get \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/rakyll/statik \
		github.com/bufbuild/buf/cmd/buf

proto:
	protoc -I . \
	-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
   --go_out ./gen/go/ --go_opt paths=source_relative \
   --go-grpc_out ./gen/go/ --go-grpc_opt paths=source_relative \
   proto/{{cookiecutter.app_name}}.proto
	protoc -I . --grpc-gateway_out ./gen/go \
	-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
     --grpc-gateway_opt logtostderr=true \
     --grpc-gateway_opt paths=source_relative \
     --grpc-gateway_opt generate_unbound_methods=true \
     proto/{{cookiecutter.app_name}}.proto
	protoc -I . -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\
		--openapiv2_out ./gen/openapiv2 --openapiv2_opt logtostderr=true proto/{{cookiecutter.app_name}}.proto

generate:
	buf generate --file ./proto/{{cookiecutter.app_name}}.proto
	# Generate static assets for OpenAPI UI
	statik -m -f -src third_party/OpenAPI/

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...


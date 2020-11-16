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
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make clean           Clean the directory tree.'
	@echo '    make dep             Update dependencies.'
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

generate:
	buf generate --file ./proto/{{cookiecutter.app_name}}.proto
	# Generate static assets for OpenAPI UI
	statik -m -f -src third_party/OpenAPI/

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
	go clean ./...

test:
	go test ./...

run: build
	./bin/${BIN_NAME}

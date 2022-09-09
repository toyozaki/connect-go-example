PACKAGE_NAME=$(shell basename $(shell pwd))
REPOSITORY_PATH=github.com/toyozaki
LDFLAGS="-s -X $(REPOSITORY_PATH)/$(PACKAGE_NAME)/cmd.appName=$(PACKAGE_NAME) -X $(REPOSITORY_PATH)/$(PACKAGE_NAME)/cmd.version=$(shell git rev-parse --short HEAD)"

.PHONY: all
all: prepare test build

.PHONY: prepare
prepare:
	go get -t -v ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -ldflags $(LDFLAGS) -o $(PACKAGE_NAME)

.PHONY: run
run:
	go run main.go

.PHONY: requirements
requirements:
	brew install clang-format
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

.PHONY: test-requirements
test-requirements:
	which buf grpcurl protoc-gen-go protoc-gen-connect-go

.PHONY: generate
generate:
	buf lint
	buf generate

.PHONY: clean
clean:
	rm $(PACKAGE_NAME)

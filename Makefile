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

.PHONY: clean
clean:
	rm $(PACKAGE_NAME)

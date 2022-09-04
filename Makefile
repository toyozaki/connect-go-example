PACKAGE_NAME=$(shell basename $(shell pwd))
LDFLAGS="-s -X main.version=$(shell git rev-parse --short HEAD)"

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

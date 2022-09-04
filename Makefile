PACKAGE_NAME=$(shell basename $(shell pwd))

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
	go build -o $(PACKAGE_NAME)

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	rm $(PACKAGE_NAME)

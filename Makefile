.PHONY: build test benchmark

BINARY_NAME=lambdamux

build:
	go build -v -o $(BINARY_NAME) .

test:
	go test -v ./...

benchmark:
	go test -bench=. -benchmem ./...

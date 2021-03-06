# Define parameters
BINARY=gamecha
SHELL := /bin/bash
GOPACKAGES = $(shell go list ./... | grep -v vendor)
ROOTDIR = $(pwd)

.PHONY: build install test linux

default: build

build: main.go config.go
	go build -v -o ./build/${BINARY} main.go config.go

install:
	go install  ./...

test:
	go test -race -cover ${GOPACKAGES}

clean:
	rm -rf build

linux: main.go config.go
	GOOS=linux GOARCH=amd64 go build -o ./build/linux/${BINARY} main.go config.go

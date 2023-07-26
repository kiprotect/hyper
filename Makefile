## simple makefile to log workflow
.PHONY: all test clean build install protobuf examples

SHELL := /bin/bash
VERSION ?= development
SD ?= default

GOFLAGS ?= $(GOFLAGS:) -ldflags "-X 'github.com/iris-connect/hyper.Version=$(VERSION)'"

export HYPER_TEST = yes

HYPER_TEST_SETTINGS ?= "$(shell pwd)/settings/test"

all: dep install

build:
	go build $(GOFLAGS) ./...

dep:
	@go get ./...

install: build
	go install $(GOFLAGS) ./...

test:
	HYPER_SETTINGS=$(HYPER_TEST_SETTINGS) go test $(testargs) `go list ./...`

test-races:
	HYPER_SETTINGS=$(HYPER_TEST_SETTINGS) go test -race $(testargs) `go list ./...`

bench:
	HYPER_SETTINGS=$(HYPER_TEST_SETTINGS) go test -run=NONE -bench=. $(GOFLAGS) `go list ./... | grep -v api/`

clean:
	@go clean $(GOFLAGS) -i ./...

copyright:
	python3 .scripts/make_copyright_headers.py

protobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protobuf/hyper.proto

certs:
	rm -rf settings/dev/certs/*
	rm -rf settings/test/certs/*
	(cd settings/dev/certs; ../../../.scripts/make_certs.sh)
	(cd settings/test/certs; ../../../.scripts/make_certs.sh)

sd-setup:
	# we always reset the directory and load the certs
	.scripts/sd_setup.sh settings/dev/directory --reset
	# then we load additional entries
	.scripts/sd_setup.sh settings/dev/directory/$(SD)

sd-test-setup:
	.scripts/sd_setup.sh settings/test/directory

examples:
	@go build $(GOFLAGS) -tags examples ./...
	@go install $(GOFLAGS) -tags examples ./...

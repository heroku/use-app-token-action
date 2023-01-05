#!make

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# Required for `run`
APP_ID ?=
PRIVATE_KEY ?=
GITHUB_REPOSITORY ?=

default: clean build

clean:
	rm -rf bin

build:
	go build -o bin/main main.go

test:
	go test $$(go list ./... | grep -v mocks)

run:
	go run main.go

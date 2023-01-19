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
	rm -rf tmp

build:
	./build.sh $(ARGS)

test:
	go test $$(go list ./... | grep -v 'cmd\|_mocks')

run:
	go run cmd/get-use-app-token-action/main.go

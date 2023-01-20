#!make

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# !IMPORTANT Set the application version here
export VERSION := 0.0.1

# Required for `run`
APP_ID ?=
PRIVATE_KEY ?=
GITHUB_REPOSITORY ?=

default: clean build

clean:
	rm -rf tmp

build:
	./ci/build.sh $(ARGS)

push-binaries:
	./ci/push-binaries.sh $(ARGS)

version:
	./ci/version.sh $(ARGS)

test:
	go test $$(go list ./... | grep -v 'cmd\|_mocks')

run:
	go run ./ $(ARGS)

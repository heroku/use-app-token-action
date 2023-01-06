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
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/main-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/main-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/main-windows-amd64 main.go
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o bin/main-windows-arm64 main.go

test:
	go test $$(go list ./... | grep -v mocks)

run:
	go run main.go

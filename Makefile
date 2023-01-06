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
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/main-macOS-X64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/main-macOS-ARM64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main-Linux-X64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/main-Linux-ARM64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/main-Windows-X64 main.go
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o bin/main-Windows-ARM64 main.go

test:
	go test $$(go list ./... | grep -v mocks)

run:
	go run main.go

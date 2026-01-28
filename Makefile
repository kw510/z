ifneq ($(filter test,$(MAKECMDGOALS)),)
	include ./.env.test
	export $(shell sed 's/=.*//' .env.test)
endif

all: build

build: $(shell find . -type f -name "*.go")
	go build -o bin/api main.go

clean:
	go clean
	if [ -f bin/api ]; then rm bin/api; fi

run:
	go run main.go

test:
	go test -race ./...

test\:ci:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
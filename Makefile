.PHONY: build, format, lint, release, run, test

build:
	go build -o ./bin/ ./cmd/ova-service-api

format:
	go fmt ./...

lint:
	golangci-lint run -v

test:
	go test -v ./...

release: format lint test build

run:
	go run ./cmd/ova-service-api

.PHONY: build, format, lint, release, run

build:
	go build -o ./bin/ ./cmd/ova-service-api

format:
	go fmt ./...

lint:
	golangci-lint run -v

release: format lint build

run:
	go run ./cmd/ova-service-api

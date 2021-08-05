.PHONY: build, format, run, lint, release

build:
	go build -o ./bin/ ./cmd/ova-service-api

run:
	go run ./cmd/ova-service-api

format:
	go fmt ./...

lint:
	golangci-lint run -v

release: format lint build

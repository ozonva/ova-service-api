.PHONY: build, format, lint, release, run, test, generate

build:
	go mod tidy
	go build -o ./bin/ ./cmd/ova-service-api

format:
	go fmt ./...

lint:
	golangci-lint run -v

test:
	go test -v ./...

generate:
	go generate ./...

release: format lint generate test build

run:
	go run ./cmd/ova-service-api

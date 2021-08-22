.PHONY: build, format, lint, release, run, test, race, clean, generate

build:
	go mod tidy
	go build -o ./bin/ ./cmd/ova-service-api

format:
	go fmt ./...

lint:
	golangci-lint run -v

test:
	go test -v ./...

race:
	go test -race ./...

clean:
	go clean -testcache

generate:
	go generate ./...

release: format lint generate test race build

run:
	go run ./cmd/ova-service-api

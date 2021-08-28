.PHONY: build, format, lint, release, run, test, race, clean, generate, deps, vendor-proto, generate-proto, proto

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

release: format lint generate proto test race build

run:
	go run ./cmd/ova-service-api

deps:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u google.golang.org/grpc
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

generate-proto:
	mkdir -p swagger
	mkdir -p pkg/ova-service-api
	protoc -I vendor.protogen \
		--go_out=pkg/ova-service-api --go_opt=paths=import \
		--go-grpc_out=pkg/ova-service-api --go-grpc_opt=paths=import \
		--grpc-gateway_out=pkg/ova-service-api \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=import \
		--swagger_out=allow_merge=true,merge_file_name=api:swagger \
		api/ova-service-api/*.proto
	mv pkg/ova-service-api/github.com/ozonva/ova-service-api/* pkg/ova-service-api/
	rm -rf pkg/ova-service-api/github.com
	mkdir -p cmd/ova-service-api
	go mod tidy

vendor-proto:
	mkdir -p vendor.protogen
	mkdir -p vendor.protogen/api/ova-service-api
	cp api/ova-service-api/service.proto vendor.protogen/api/ova-service-api/service.proto
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi

proto: vendor-proto generate-proto

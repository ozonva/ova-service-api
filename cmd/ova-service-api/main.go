package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozonva/ova-service-api/internal/api"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

const (
	httpServerEndpoint = "localhost:8081"
	grpcServerEndpoint = "localhost:8082"
)

func main() {
	go runHttpServer()

	if err := runGrpcServer(); err != nil {
		log.Fatal(err)
	}
}

func runGrpcServer() error {
	listen, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		log.Fatalf("gRPC: failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServiceAPIServer(s, api.NewGrpcApiServer())

	if grpcErr := s.Serve(listen); grpcErr != nil {
		log.Fatalf("gRPC: failed to serve: %v", grpcErr)
	}

	return nil
}

func runHttpServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := pb.RegisterServiceAPIHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts); err != nil {
		log.Fatalf("http: failed to register server: %v", err)
	}

	if httpErr := http.ListenAndServe(httpServerEndpoint, mux); httpErr != nil {
		log.Fatalf("http: failed to listen: %v", httpErr)
	}
}

package main

import (
	"github.com/ozonva/ova-service-api/internal/api"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

const grpcServerEndpoint = "localhost:8082"

func main() {
	if err := runGrpcServer(); err != nil {
		panic(err)
	}
}

func runGrpcServer() error {
	listen, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServiceAPIServer(s, api.NewGrpcApiServer())

	if grpcErr := s.Serve(listen); grpcErr != nil {
		log.Fatalf("failed to serve: %v", grpcErr)
	}

	return nil
}

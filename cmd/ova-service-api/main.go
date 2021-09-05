package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/ozonva/ova-service-api/internal/api"
	"github.com/ozonva/ova-service-api/internal/repo"

	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

const (
	httpServerEndpoint = "localhost:8081"
	grpcServerEndpoint = "localhost:8082"
)

func main() {
	// Ignore error because .env file may not exist, in this case real environment variables will be used
	_ = godotenv.Load()

	dsn, ok := os.LookupEnv("DATABASE_CONNECTION_STRING")

	if !ok {
		log.Fatal("DATABASE_CONNECTION_STRING environment variable is required")
	}

	ctx := context.Background()
	pgRepo, err := repo.NewPostgresServiceRepo(ctx, dsn)

	if err != nil {
		log.Fatal(err)
	}

	go runHttpServer(ctx)

	if err = runGrpcServer(ctx, pgRepo); err != nil {
		log.Fatal(err)
	}
}

// Actually it should use root context, but for this task we do not use it
func runGrpcServer(_ context.Context, repo repo.Repo) error {
	listen, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		log.Fatalf("gRPC: failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServiceAPIServer(s, api.NewGrpcApiServer(repo))

	if grpcErr := s.Serve(listen); grpcErr != nil {
		log.Fatalf("gRPC: failed to serve: %v", grpcErr)
	}

	return nil
}

func runHttpServer(ctx context.Context) {
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

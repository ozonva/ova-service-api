package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/ozonva/ova-service-api/internal/api"
	flusher_ "github.com/ozonva/ova-service-api/internal/flusher"
	"github.com/ozonva/ova-service-api/internal/infrastructure/kafka"
	"github.com/ozonva/ova-service-api/internal/infrastructure/metrics"
	repo_ "github.com/ozonva/ova-service-api/internal/repo"
	saver_ "github.com/ozonva/ova-service-api/internal/saver"
	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

const (
	httpServerEndpoint   = "localhost:8081"
	grpcServerEndpoint   = "localhost:8082"
	metricEndpoint       = "localhost:9100"
	multiCreateBatchSize = 5
	flushTimeout         = 1 * time.Second
	localCapacity        = 10
	kafkaTopic           = "services"
)

func main() {
	ctx := context.Background()
	env, err := readEnvironment()
	if err != nil {
		log.Fatalf("Error occured during read environment: %s", err.Error())
	}

	resolver := newDependencyResolver(ctx, env)
	deps, err := resolver.resolve()
	if err != nil {
		log.Fatalf("Error occured during dependency resolve: %s", err.Error())
	}
	defer resolver.close()

	go runMetricServer()
	go runHttpServer(ctx)

	if err = runGrpcServer(ctx, deps.Repo, deps.Saver, deps.Flusher, deps.Producer, deps.Metrics); err != nil {
		log.Fatal(err)
	}
}

// Actually it should use root context, but for this task we do not use it
func runGrpcServer(_ context.Context, repo repo_.Repo, saver saver_.Saver, flusher flusher_.Flusher, producer kafka.Producer, metrics metrics.Metrics) error {
	listen, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		log.Fatalf("gRPC: failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor))
	pb.RegisterServiceAPIServer(server, api.NewGrpcApiServer(repo, saver, flusher, producer, metrics))

	if grpcErr := server.Serve(listen); grpcErr != nil {
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

func runMetricServer() {
	http.Handle("/metrics", promhttp.Handler())
	if httpErr := http.ListenAndServe(metricEndpoint, nil); httpErr != nil {
		log.Fatalf("metric: failed to listen: %v", httpErr)
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	flusher_ "github.com/ozonva/ova-service-api/internal/flusher"
	repo_ "github.com/ozonva/ova-service-api/internal/repo"
	saver_ "github.com/ozonva/ova-service-api/internal/saver"
)

const (
	multiCreateBatchSize = 5
	flushTimeout         = 1 * time.Second
	localCapacity        = 10
)

type Dependencies struct {
	ctx     context.Context
	repo    repo_.Repo
	flusher flusher_.Flusher
	saver   saver_.Saver
}

func NewDependencies() (*Dependencies, error) {
	return resolve()
}

func (deps *Dependencies) Close() {
	if deps.saver != nil {
		deps.saver.Close()
	}
}

func resolve() (*Dependencies, error) {
	ctx := context.Background()

	// Ignore error because .env file may not exist, in this case real environment variables will be used
	_ = godotenv.Load()

	dsn, ok := os.LookupEnv("DATABASE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("DATABASE_CONNECTION_STRING environment variable is required")
	}

	pgRepo, err := repo_.NewPostgresServiceRepo(ctx, dsn)
	if err != nil {
		return nil, err
	}

	flusher := flusher_.New(multiCreateBatchSize, pgRepo)
	saver := saver_.New(localCapacity, flushTimeout, flusher)
	saver.Init()

	deps := Dependencies{
		ctx:     ctx,
		repo:    pgRepo,
		flusher: flusher,
		saver:   saver,
	}

	return &deps, nil
}

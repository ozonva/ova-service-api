package main

import (
	"context"

	flusher_ "github.com/ozonva/ova-service-api/internal/flusher"
	"github.com/ozonva/ova-service-api/internal/kafka"
	repo_ "github.com/ozonva/ova-service-api/internal/repo"
	saver_ "github.com/ozonva/ova-service-api/internal/saver"
)

type dependencies struct {
	Repo     repo_.Repo
	Flusher  flusher_.Flusher
	Saver    saver_.Saver
	Producer kafka.Producer
}

type dependencyResolver struct {
	ctx  context.Context
	env  environment
	deps *dependencies
}

func newDependencyResolver(ctx context.Context, env environment) dependencyResolver {
	return dependencyResolver{
		ctx: ctx,
		env: env,
	}
}

func (dr *dependencyResolver) resolve() (*dependencies, error) {
	pgRepo, err := repo_.NewPostgresServiceRepo(dr.ctx, dr.env.DSN)
	if err != nil {
		return nil, err
	}

	flusher := flusher_.New(multiCreateBatchSize, pgRepo)
	saver := saver_.New(localCapacity, flushTimeout, flusher)
	saver.Init()

	producer, err := kafka.NewSyncProducer(kafkaTopic, dr.env.Brokers)
	if err != nil {
		return nil, err
	}

	deps := dependencies{
		Repo:     pgRepo,
		Flusher:  flusher,
		Saver:    saver,
		Producer: producer,
	}

	return &deps, nil
}

func (dr *dependencyResolver) close() {
	if dr.deps != nil && dr.deps.Saver != nil {
		dr.deps.Saver.Close()
	}
}

package repo

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ozonva/ova-service-api/internal/models"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresServiceRepo struct {
	ctx context.Context
	db  *sql.DB
}

func NewPostgresServiceRepo(ctx context.Context, dsn string) (*PostgresServiceRepo, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Err(err).Msg("Can't load pgx driver")
		return nil, err
	}

	if connErr := db.PingContext(ctx); connErr != nil {
		log.Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	return &PostgresServiceRepo{
		ctx: ctx,
		db:  db,
	}, nil
}

func (repo *PostgresServiceRepo) AddServices(services []models.Service) error {
	log.Debug().Msg("PostgresServiceRepo.AddServices call")

	if len(services) == 0 {
		return nil
	}

	tx, err := repo.db.BeginTx(repo.ctx, nil)
	defer func(tx *sql.Tx) {
		// Intentionally ignore the error, because tx.commit could be called before deferred rollback
		_ = tx.Rollback()
	}(tx)

	if err != nil {
		log.Err(err).Msg("Failed to begin transaction")
		return err
	}

	query := `INSERT INTO services (id, user_id, description, service_name, service_address, when_local, when_utc)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, service := range services {
		if _, err = tx.ExecContext(repo.ctx, query, service.ID, service.UserID, service.Description, service.ServiceName, service.ServiceAddress, service.WhenLocal, service.WhenUTC); err != nil {
			log.Err(err).Msg("Failed to begin transaction")
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Err(err).Msg("Error occurred during commit transaction")
		return err
	}

	log.Info().Msg("Services was successfully stored in the database")
	return nil
}

func (repo *PostgresServiceRepo) ListServices(limit uint64, offset uint64) ([]models.Service, error) {
	log.Debug().Msg("PostgresServiceRepo.ListServices call")

	panic("implement me")
}

func (repo *PostgresServiceRepo) DescribeService(serviceID uuid.UUID) (*models.Service, error) {
	log.Debug().Msg("PostgresServiceRepo.DescribeService call")

	panic("implement me")
}

func (repo *PostgresServiceRepo) RemoveService(serviceID uuid.UUID) error {
	log.Debug().Msg("PostgresServiceRepo.RemoveService call")

	panic("implement me")
}

package repo

import (
	"context"
	"database/sql"
	"fmt"

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
		txErr := tx.Rollback()
		if txErr != sql.ErrTxDone {
			log.Err(txErr).Msg("Can' rollback the transaction")
		}
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

	var (
		query string
		rows  *sql.Rows
		err   error
	)

	// This is actually a hack to handle the difference between the required Repo API which includes limit and offset
	// and gRPC server API which allows to list all.
	if limit < ^uint64(0) {
		query = `SELECT id, user_id, description, service_name, service_address, when_local, when_utc
			FROM services
			ORDER BY when_utc DESC
			LIMIT $1 OFFSET $2`
		rows, err = repo.db.QueryContext(repo.ctx, query, limit, offset)
	} else {
		query = `SELECT id, user_id, description, service_name, service_address, when_local, when_utc
			FROM services
			ORDER BY when_utc DESC`
		rows, err = repo.db.QueryContext(repo.ctx, query)
	}

	if err != nil {
		log.Err(err).Msg("Error occurred during query execution")
		return nil, err
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			log.Err(err).Msg("Can't properly close rows cursor")
		}
	}(rows)

	services := make([]models.Service, 0)

	for rows.Next() {
		var service models.Service
		if err = rows.Scan(&service.ID, &service.UserID, &service.Description, &service.ServiceName,
			&service.ServiceAddress, &service.WhenLocal, &service.WhenUTC); err != nil {
			log.Err(err).Msg("Can't parse single row")
			return nil, err
		}
		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		log.Err(err).Msg("Error occurs during cursor iteration")
		return nil, err
	}

	return services, nil
}

func (repo *PostgresServiceRepo) DescribeService(serviceID uuid.UUID) (*models.Service, error) {
	log.Debug().Msg("PostgresServiceRepo.DescribeService call")

	query := `SELECT id, user_id, description, service_name, service_address, when_local, when_utc
			FROM services
			WHERE id = $1`

	row := repo.db.QueryRowContext(repo.ctx, query, serviceID)

	var service *models.Service
	err := row.Scan(&service.ID, &service.UserID, &service.Description, &service.ServiceName,
		&service.ServiceAddress, &service.WhenLocal, &service.WhenUTC)

	switch err {
	case nil:
		return service, nil
	case sql.ErrNoRows:
		notFoundErr := fmt.Errorf("service with ID: %s was not found in the repo", serviceID.String())
		log.Err(notFoundErr).Msg("Error occurred during describe service")
		return nil, notFoundErr
	default:
		log.Err(err).Msg("Error occurred during query execution")
		return nil, err
	}
}

func (repo *PostgresServiceRepo) RemoveService(serviceID uuid.UUID) error {
	log.Debug().Msg("PostgresServiceRepo.RemoveService call")

	query := `DELETE
			FROM services
			WHERE id = $1`

	if _, err := repo.db.ExecContext(repo.ctx, query, serviceID); err != nil {
		log.Err(err).Msg("Error occurs during delete operation execution")
		return err
	}

	return nil
}

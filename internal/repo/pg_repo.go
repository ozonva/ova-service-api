package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ozonva/ova-service-api/internal/models"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type dbService struct {
	ID             uuid.UUID
	UserID         uint64
	Description    sql.NullString
	ServiceName    sql.NullString
	ServiceAddress sql.NullString
	WhenLocal      sql.NullTime
	WhenUTC        sql.NullTime
}

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

	queryParts := make([]string, len(services)+1)
	queryParts[0] = "INSERT INTO services (id, user_id, description, service_name, service_address, when_local, when_utc) VALUES "

	queryValues := make([]interface{}, 0)

	for i, service := range services {
		// This looks crazy, but this is because of placeholders in pgx: $1, $2, $3, etc. I can't find a way how to do it easier using database/sql.
		queryParts[i+1] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", 7*i+1, 7*i+2, 7*i+3, 7*i+4, 7*i+5, 7*i+6, 7*i+7)
		queryValues = append(queryValues, service.ID, service.UserID, service.Description, service.ServiceName, service.ServiceAddress, service.WhenLocal, service.WhenUTC)
	}

	query := strings.Join(queryParts, "")
	query = query[:len(query)-1] // Remove trailing comma

	if _, err := repo.db.ExecContext(repo.ctx, query, queryValues...); err != nil {
		log.Err(err).Msg("Failed to begin transaction")
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
		var service dbService

		if err = rows.Scan(&service.ID, &service.UserID, &service.Description, &service.ServiceName,
			&service.ServiceAddress, &service.WhenLocal, &service.WhenUTC); err != nil {
			log.Err(err).Msg("Can't parse single row")
			return nil, err
		}

		services = append(services, mapDBServiceToDomainService(&service))
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

	var service dbService
	err := row.Scan(&service.ID, &service.UserID, &service.Description, &service.ServiceName,
		&service.ServiceAddress, &service.WhenLocal, &service.WhenUTC)

	switch err {
	case nil:
		domainService := mapDBServiceToDomainService(&service)
		return &domainService, nil
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

func (repo *PostgresServiceRepo) UpdateService(service *models.Service) error {
	log.Debug().Msg("PostgresServiceRepo.UpdateService call")

	if service == nil {
		nilErr := fmt.Errorf("service is nil")
		log.Err(nilErr).Msg("Error occurred during update service")
		return nilErr
	}

	query := `SELECT version
			FROM services
			WHERE id = $1`

	row := repo.db.QueryRowContext(repo.ctx, query, service.ID)
	var version int

	err := row.Scan(&version)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		notFoundErr := fmt.Errorf("service with ID: %s was not found in the repo", service.ID.String())
		log.Err(notFoundErr).Msg("Error occurred during update service")
		return notFoundErr
	default:
		log.Err(err).Msg("Error occurred during query execution")
		return err
	}

	query = `UPDATE services
			SET user_id = $1,
			    description = $2,
			    service_name = $3,
			    service_address = $4,
			    when_local = $5,
			    when_utc = $6,
			    updated_at = CURRENT_TIMESTAMP,
			    version = version + 1
			WHERE id = $7
			  AND version = $8`

	res, err := repo.db.ExecContext(repo.ctx, query, service.UserID, service.Description, service.ServiceName,
		service.ServiceAddress, service.WhenLocal, service.WhenUTC, service.ID, version)

	if err != nil {
		log.Err(err).Msg("Error occurs during update operation execution")
		return err
	}

	cnt, err := res.RowsAffected()

	if err != nil {
		log.Err(err).Msg("Error occurs during update operation execution")
		return err
	}

	if cnt == 0 {
		concurrencyErr := fmt.Errorf("service with ID: %s was not updated because entity already changed by other request", service.ID.String())
		log.Err(concurrencyErr).Msg("Error occurs during update operation execution")
		return concurrencyErr
	}

	log.Info().Msg("Service was successfully updated")
	return nil
}

func mapDBServiceToDomainService(service *dbService) models.Service {
	var domainService models.Service

	domainService.ID = service.ID
	domainService.UserID = service.UserID

	// We can skip Valid check because default string value is OK for us
	domainService.Description = service.Description.String
	domainService.ServiceName = service.ServiceName.String
	domainService.ServiceAddress = service.ServiceAddress.String

	if service.WhenLocal.Valid {
		domainService.WhenLocal = &service.WhenLocal.Time
	}
	if service.WhenUTC.Valid {
		domainService.WhenUTC = &service.WhenUTC.Time
	}

	return domainService
}

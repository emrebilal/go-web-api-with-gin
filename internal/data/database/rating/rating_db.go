package rating

import (
	"context"
	"database/sql"
	"errors"
	"rating-api/internal/util/env"
	"rating-api/internal/util/logger"
	"rating-api/internal/util/validator"
	"time"

	_ "github.com/lib/pq"
)

type IRatingDb interface {
	AddRate(ch chan *AddRatingResponse, model *AddRatingModel)
	GetAllRate(ch chan *GetAllRatingsResponse, model *GetAllRatingsModel)
}

type RatingDb struct {
	loggr            logger.ILogger
	validatr         validator.IValidator
	environment      env.IEnvironment
	connectionString string
	driverName       string
	timeout          time.Duration
}

// NewRatingDb
// Returns a new RatingDb.
func NewRatingDb(loggr logger.ILogger, validatr validator.IValidator, environment env.IEnvironment) IRatingDb {
	db := RatingDb{
		environment:      environment,
		loggr:            loggr,
		validatr:         validatr,
		driverName:       "postgres",
		connectionString: environment.Get(env.PostgresqlConnectionString),
		timeout:          time.Second * 5,
	}

	return &db
}

// AddRate
// Add rating for a service provider.
func (d *RatingDb) AddRate(ch chan *AddRatingResponse, model *AddRatingModel) {
	modelErr := d.validatr.ValidateStruct(model)
	if modelErr != nil {
		d.loggr.Error(modelErr.Error())
		ch <- &AddRatingResponse{Error: modelErr}
		return
	}

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		d.loggr.Error(err.Error())
		ch <- &AddRatingResponse{Error: err}
		return
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	query := `insert into ratings (username, provider_id, service_id, rate, created_date) 
				values ($1, $2, $3, $4, current_timestamp)
				on conflict(service_id)
				do nothing`

	result, dbErr := connection.ExecContext(ctx, query, model.UserName, model.ProviderId, model.ServiceId, model.Rate)
	if dbErr != nil {
		d.loggr.Error(dbErr.Error())
		ch <- &AddRatingResponse{Error: dbErr}
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		d.loggr.Error(err.Error())
		ch <- &AddRatingResponse{Error: err}
		return
	}

	if rows != 1 {
		d.loggr.Error("could not add rate")
		ch <- &AddRatingResponse{Error: errors.New("could not add rate")}
		return
	}

	ch <- &AddRatingResponse{}
}

// GetAllRate
// Get all ratings for a service provider.
func (d *RatingDb) GetAllRate(ch chan *GetAllRatingsResponse, model *GetAllRatingsModel) {
	modelErr := d.validatr.ValidateStruct(model)
	if modelErr != nil {
		d.loggr.Error(modelErr.Error())
		ch <- &GetAllRatingsResponse{Error: modelErr}
		return
	}

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		d.loggr.Error(err.Error())
		ch <- &GetAllRatingsResponse{Error: err}
		return
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	query := `select rate from ratings where provider_id = $1`

	rows, dbErr := connection.QueryContext(ctx, query, model.ProviderId)

	if dbErr != nil && dbErr != sql.ErrNoRows {
		d.loggr.Error(dbErr.Error())
		ch <- &GetAllRatingsResponse{Error: dbErr}
		return
	}
	if dbErr == sql.ErrNoRows {
		ch <- &GetAllRatingsResponse{Rates: []int{}}
		return
	}

	var response GetAllRatingsResponse
	for rows.Next() {
		var rate int
		if err := rows.Scan(&rate); err != nil {
			d.loggr.Error(err.Error())
			ch <- &GetAllRatingsResponse{Error: err}
			return
		}
		response.Rates = append(response.Rates, rate)
	}

	ch <- &response
}

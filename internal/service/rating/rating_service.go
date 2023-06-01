package rating

import (
	"errors"
	"rating-api/internal/data/database/rating"
	"rating-api/internal/util/env"
	"rating-api/internal/util/logger"
	"rating-api/internal/util/validator"
)

type IRatingService interface {
	SendRating(ch chan *SendRatingServiceResponse, model *SendRatingServiceModel)
	GetAverageRating(ch chan *GetAverageRatingServiceResponse, model *GetAverageRatingServiceModel)
}

type RatingService struct {
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	ratingDb    rating.IRatingDb
}

// NewRatingService
// Returns a new RatingService.
func NewRatingService(
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	ratingDb rating.IRatingDb,
) IRatingService {
	service := RatingService{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
	}

	if ratingDb != nil {
		service.ratingDb = ratingDb
	} else {
		service.ratingDb = rating.NewRatingDb(loggr, validatr, environment)
	}

	return &service
}

func (r *RatingService) SendRating(ch chan *SendRatingServiceResponse, model *SendRatingServiceModel) {
	modelErr := r.validatr.ValidateStruct(model)
	if modelErr != nil {
		r.loggr.Error(modelErr.Error())
		ch <- &SendRatingServiceResponse{Error: modelErr}
		return
	}

	chRatingDb := make(chan *rating.AddRatingResponse)
	defer close(chRatingDb)

	go r.ratingDb.AddRate(chRatingDb, &rating.AddRatingModel{
		UserName:   model.UserName,
		ProviderId: model.ProviderId,
		ServiceId:  model.ServiceId,
		Rate:       model.Rate,
	})

	dbResponse := <-chRatingDb
	if dbResponse.Error != nil {
		ch <- &SendRatingServiceResponse{Error: dbResponse.Error}
		return
	}

	ch <- &SendRatingServiceResponse{Info: "Added rating for ServiceId: " + model.ServiceId + " getting from ProviderId: " + model.ProviderId}
}

func (r *RatingService) GetAverageRating(ch chan *GetAverageRatingServiceResponse, model *GetAverageRatingServiceModel) {
	modelErr := r.validatr.ValidateStruct(model)
	if modelErr != nil {
		r.loggr.Error(modelErr.Error())
		ch <- &GetAverageRatingServiceResponse{Error: modelErr}
		return
	}

	chRatingDb := make(chan *rating.GetAllRatingsResponse)
	defer close(chRatingDb)

	go r.ratingDb.GetAllRate(chRatingDb, &rating.GetAllRatingsModel{
		ProviderId: model.ProviderId,
	})

	dbResponse := <-chRatingDb
	if dbResponse.Error != nil {
		ch <- &GetAverageRatingServiceResponse{Error: dbResponse.Error}
		return
	}

	if len(dbResponse.Rates) == 0 {
		r.loggr.Error("No ratings found for ProviderId: " + model.ProviderId)
		ch <- &GetAverageRatingServiceResponse{Error: errors.New("No ratings found for ProviderId: " + model.ProviderId)}
		return
	}

	// calculate average rate
	sum := 0
	for i := 0; i < len(dbResponse.Rates); i++ {
		sum += (dbResponse.Rates[i])
	}

	avg := (float64(sum)) / (float64(len(dbResponse.Rates)))

	ch <- &GetAverageRatingServiceResponse{AverageRating: AverageRatingModel{ProviderId: model.ProviderId, AverageRate: avg}}
}

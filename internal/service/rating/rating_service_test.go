package rating

import (
	"errors"
	ratingDb "rating-api/internal/data/database/rating"
	"rating-api/internal/util/env"
	"rating-api/internal/util/logger"
	"rating-api/internal/util/validator"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RatingServiceTestSuite struct {
	suite.Suite
	ratingService   IRatingService
	mockEnvironment *env.MockIEnvironment
	mockLogger      *logger.MockILogger
	mockValidator   *validator.MockIValidator
	mockRatingDb    *ratingDb.MockIRatingDb
}

// Run suite.
func TestService(t *testing.T) {
	suite.Run(t, new(RatingServiceTestSuite))
}

// Runs before each test in the suite.
func (r *RatingServiceTestSuite) SetupTest() {
	r.T().Log("Setup")

	ctrl := gomock.NewController(r.T())
	defer ctrl.Finish()

	r.mockEnvironment = env.NewMockIEnvironment(ctrl)
	r.mockLogger = logger.NewMockILogger(ctrl)
	r.mockValidator = validator.NewMockIValidator(ctrl)
	r.mockRatingDb = ratingDb.NewMockIRatingDb(ctrl)

	r.ratingService = NewRatingService(r.mockEnvironment, r.mockLogger, r.mockValidator, r.mockRatingDb)
}

// Runs after each test in the suite.
func (r *RatingServiceTestSuite) TearDownTest() {
	r.T().Log("Teardown")
}

func (r *RatingServiceTestSuite) TestSendRating_HappyPath_Success() {
	model := SendRatingServiceModel{
		UserName:   "emre.bilal",
		ProviderId: "test-1",
		ServiceId:  "s-1",
		Rate:       4,
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(nil)

	r.mockRatingDb.
		EXPECT().
		AddRate(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ch chan *ratingDb.AddRatingResponse, model *ratingDb.AddRatingModel) {
				ch <- &ratingDb.AddRatingResponse{
					Error: nil,
				}
			},
		)

	ch := make(chan *SendRatingServiceResponse)
	defer close(ch)

	go r.ratingService.SendRating(ch, &model)
	response := <-ch

	r.Nil(response.Error)
}

func (r *RatingServiceTestSuite) TestSendRating_ModelValidationError_ReturnsError() {
	model := SendRatingServiceModel{
		UserName:   "emre.bilal",
		ProviderId: "test-1",
		ServiceId:  "s-1",
		Rate:       0,
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(errors.New("Rate must be between 1 and 5"))

	r.mockLogger.EXPECT().Error(gomock.Any())

	ch := make(chan *SendRatingServiceResponse)
	defer close(ch)

	go r.ratingService.SendRating(ch, &model)
	response := <-ch

	r.NotNil(response.Error)
	r.EqualError(response.Error, "Rate must be between 1 and 5")
}

func (r *RatingServiceTestSuite) TestSendRating_DatabaseError_ReturnsError() {
	model := SendRatingServiceModel{
		UserName:   "emre.bilal",
		ProviderId: "test-1",
		ServiceId:  "s-1",
		Rate:       4,
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(nil)

	r.mockRatingDb.
		EXPECT().
		AddRate(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ch chan *ratingDb.AddRatingResponse, model *ratingDb.AddRatingModel) {
				ch <- &ratingDb.AddRatingResponse{
					Error: errors.New("an error occurred"),
				}
			},
		)

	ch := make(chan *SendRatingServiceResponse)
	defer close(ch)

	go r.ratingService.SendRating(ch, &model)
	response := <-ch

	r.NotNil(response.Error)
	r.Error(response.Error)
}

func (r *RatingServiceTestSuite) TestGetAverageRating_HappyPath_Success() {
	model := GetAverageRatingServiceModel{
		ProviderId: "test-1",
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(nil)

	r.mockRatingDb.
		EXPECT().
		GetAllRate(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ch chan *ratingDb.GetAllRatingsResponse, model *ratingDb.GetAllRatingsModel) {
				ch <- &ratingDb.GetAllRatingsResponse{
					Error: nil,
					Rates: []int{4, 5, 4, 3},
				}
			},
		)

	avgRate := (4 + 5 + 4 + 3) / 4

	ch := make(chan *GetAverageRatingServiceResponse)
	defer close(ch)

	go r.ratingService.GetAverageRating(ch, &model)
	response := <-ch

	r.Nil(response.Error)
	r.Equal(response.AverageRating.AverageRate, float64(avgRate))
}

func (r *RatingServiceTestSuite) TestGetAverageRating_ModelValidationError_ReturnsError() {
	model := GetAverageRatingServiceModel{
		ProviderId: "",
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(errors.New("ProviderId cannot be empty"))

	r.mockLogger.EXPECT().Error(gomock.Any())

	ch := make(chan *GetAverageRatingServiceResponse)
	defer close(ch)

	go r.ratingService.GetAverageRating(ch, &model)
	response := <-ch

	r.NotNil(response.Error)
	r.EqualError(response.Error, "ProviderId cannot be empty")
}

func (r *RatingServiceTestSuite) TestGetAverageRating_DatabaseError_ReturnsError() {
	model := GetAverageRatingServiceModel{
		ProviderId: "test-1",
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(nil)

	r.mockRatingDb.
		EXPECT().
		GetAllRate(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ch chan *ratingDb.GetAllRatingsResponse, model *ratingDb.GetAllRatingsModel) {
				ch <- &ratingDb.GetAllRatingsResponse{
					Error: errors.New("an error occurred"),
					Rates: []int{},
				}
			},
		)

	ch := make(chan *GetAverageRatingServiceResponse)
	defer close(ch)

	go r.ratingService.GetAverageRating(ch, &model)
	response := <-ch

	r.NotNil(response.Error)
	r.Error(response.Error)
}

func (r *RatingServiceTestSuite) TestGetAverageRating_NoRatingsFound_ReturnsError() {
	model := GetAverageRatingServiceModel{
		ProviderId: "test-1",
	}

	r.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&model)).
		Return(nil)

	r.mockRatingDb.
		EXPECT().
		GetAllRate(gomock.Any(), gomock.Any()).
		DoAndReturn(
			func(ch chan *ratingDb.GetAllRatingsResponse, model *ratingDb.GetAllRatingsModel) {
				ch <- &ratingDb.GetAllRatingsResponse{
					Error: nil,
					Rates: []int{},
				}
			},
		)

	r.mockLogger.EXPECT().Error(gomock.Any())

	ch := make(chan *GetAverageRatingServiceResponse)
	defer close(ch)

	go r.ratingService.GetAverageRating(ch, &model)
	response := <-ch

	r.NotNil(response.Error)
	r.EqualError(response.Error, "No ratings found for ProviderId: "+model.ProviderId)
}

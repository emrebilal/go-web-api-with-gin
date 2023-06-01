package rating

import (
	"net/http"
	"rating-api/internal/api"
	"rating-api/internal/service/rating"
	"rating-api/internal/util/env"
	"rating-api/internal/util/logger"
	"rating-api/internal/util/validator"

	"github.com/gin-gonic/gin"
)

type IRatingController interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
	AddRating(context *gin.Context)
	GetAverageRating(context *gin.Context)
}

type RatingController struct {
	path          string
	environment   env.IEnvironment
	loggr         logger.ILogger
	validatr      validator.IValidator
	ratingService rating.IRatingService
}

// NewRatingController
// Returns a new RatingController.
func NewRatingController(
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	ratingService rating.IRatingService,
) IRatingController {
	controller := RatingController{
		path:        "rating",
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
	}

	if ratingService != nil {
		controller.ratingService = ratingService
	} else {
		controller.ratingService = rating.NewRatingService(environment, loggr, validatr, nil)
	}

	return &controller
}

// RegisterRoutes
// Registers routes to gin.
func (c *RatingController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routes := routerGroup.Group(c.path)
	routes.POST("add", c.AddRating)
	routes.GET("avg", c.GetAverageRating)
}

// AddRating
//
//	@basePath		/api
//	@router			/v1/rating/add [post]
//	@tags			Rating
//	@summary		Add provider rating.
//	@description	Add provider rating.
//	@accept			json
//	@produce		json
//	@success		200		{object}	api.ApiResponse
//	@failure		400		{object}	api.ApiResponse
//	@failure		401		{object}	api.ApiResponse
//	@failure		500		{object}	api.ApiResponse
//
//	@Param			Model	body		AddRatingModel	true	"Request model"
func (c *RatingController) AddRating(context *gin.Context) {
	var model AddRatingModel
	err := context.ShouldBindJSON(&model)
	if err != nil {
		context.Error(err)
		context.JSON(http.StatusBadRequest, api.RespondError(err.Error()))
		return
	}

	chRatingService := make(chan *rating.SendRatingServiceResponse)
	defer close(chRatingService)

	go c.ratingService.SendRating(chRatingService, &rating.SendRatingServiceModel{
		UserName:   model.UserName,
		ProviderId: model.ProviderId,
		ServiceId:  model.ServiceId,
		Rate:       model.Rate,
	})

	ratingServiceResponse := <-chRatingService
	if ratingServiceResponse.Error != nil {
		context.Error(ratingServiceResponse.Error)
		context.JSON(http.StatusBadRequest, api.RespondError(ratingServiceResponse.Error.Error()))
		return
	}

	context.JSON(http.StatusOK, api.RespondOk(ratingServiceResponse))
}

// GetAverageRating
//
//	@basePath		/api
//	@router			/v1/rating/avg [get]
//	@tags			Rating
//	@summary		Get provider's average rating.
//	@description	Get provider's average rating.
//	@accept			json
//	@produce		json
//	@success		200			{object}	api.ApiResponse
//	@failure		400			{object}	api.ApiResponse
//	@failure		401			{object}	api.ApiResponse
//	@failure		500			{object}	api.ApiResponse
//	@Param			providerId	query		string	true	"Provider Id"
func (c *RatingController) GetAverageRating(context *gin.Context) {
	providerId := context.Query("providerId")

	chRatingService := make(chan *rating.GetAverageRatingServiceResponse)
	defer close(chRatingService)

	go c.ratingService.GetAverageRating(chRatingService, &rating.GetAverageRatingServiceModel{
		ProviderId: providerId,
	})

	ratingServiceResponse := <-chRatingService
	if ratingServiceResponse.Error != nil {
		context.Error(ratingServiceResponse.Error)
		context.JSON(http.StatusBadRequest, api.RespondError(ratingServiceResponse.Error.Error()))
		return
	}

	context.JSON(http.StatusOK, api.RespondOk(ratingServiceResponse))
}

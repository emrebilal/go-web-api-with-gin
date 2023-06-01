package main

import (
	"fmt"
	"rating-api/docs"
	"rating-api/internal/api"
	"rating-api/internal/api/controller/v1/health"
	"rating-api/internal/api/controller/v1/rating"
	"rating-api/internal/util/env"
	"rating-api/internal/util/logger"
	"rating-api/internal/util/validator"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//	@title			Rating API
//	@version		1.0
//	@description	This is an rating service for providers.
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8080
//	@BasePath		/api
//	@accept			json
//	@produce		json
//	@schemes		http https
func main() {
	environment := env.New()
	environment.Init()
	loggr := logger.New(environment)
	defer loggr.Sync()
	validatr := validator.New()

	router := gin.New()
	router.Use(api.LoggingMiddleware(loggr))
	addRoutes(router, environment, loggr, validatr)
	addSwagger(router, environment)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	router.Run()
}

func addRoutes(router *gin.Engine, environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator) {
	api := router.Group("api")
	health.NewHealthController().RegisterRoutes(api)

	v1 := api.Group("v1")
	rating.NewRatingController(environment, loggr, validatr, nil).RegisterRoutes(v1)
}

func addSwagger(router *gin.Engine, environment env.IEnvironment) {
	docs.SwaggerInfo.Title = fmt.Sprintf("Rating API (%v)", environment.Get(env.AppEnvironment))
	docs.SwaggerInfo.Host = environment.Get(env.AppHost)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

package health

import (
	"net/http"
	"rating-api/internal/api"

	"github.com/gin-gonic/gin"
)

type IHealthController interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
	Ping(context *gin.Context)
}

type HealthController struct{}

// NewHealthController
// Returns a new HealthController.
func NewHealthController() IHealthController {
	return &HealthController{}
}

func (c *HealthController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("ping", c.Ping)
}

// Ping
//
//	@basePath		/api
//	@router			/ping [get]
//	@tags			Health
//	@summary		Send a ping request.
//	@description	Send a ping request.
//	@accept			json
//	@produce		json
//	@success		200	{object}	api.ApiResponse
//	@failure		400	{object}	api.ApiResponse
//	@failure		500	{object}	api.ApiResponse
func (c *HealthController) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, api.RespondOk("Ping OK"))
}

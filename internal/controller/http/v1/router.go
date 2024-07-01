package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"warehousesAPI/internal/usecase"
	"warehousesAPI/pkg/logger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(
	handler *gin.Engine,
	l logger.Interface,
	i usecase.ItemsUseCase,
	r usecase.ReservationsUsecase,
	w usecase.WarehousesUsecase,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newWarehousesAPIRoutes(h, w, l)
		newReservationsAPIRoutes(h, r, l)
		newItemsAPIRoutes(h, i, l)
	}
}

package v1

import (
	"github.com/gin-gonic/gin"
	"warehousesAPI/internal/entity"
	"warehousesAPI/internal/usecase"
	"warehousesAPI/pkg/logger"
)

type warehousesAPIRoutes struct {
	warehouses usecase.WarehousesUsecase
	l          logger.Interface
}

func newWarehousesAPIRoutes(handler *gin.RouterGroup, w usecase.WarehousesUsecase, l logger.Interface) {
	r := &warehousesAPIRoutes{w, l}

	h := handler.Group("/warehouses")
	{
		h.POST("/", r.createWarehouse)
	}

}

func (w *warehousesAPIRoutes) createWarehouse(c *gin.Context) {
	var wh entity.Warehouse
	err := c.ShouldBindJSON(&wh)
	if err != nil {
		return
	}

	if err := w.warehouses.WarehouseCreate(c.Request.Context(), wh); err != nil {
		return
	}
}

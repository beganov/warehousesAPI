package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"warehousesAPI/internal/entity"
	"warehousesAPI/internal/usecase"
	"warehousesAPI/pkg/logger"
)

type itemsAPIRouter struct {
	items usecase.ItemsUseCase
	l     logger.Interface
}

func newItemsAPIRoutes(handler *gin.RouterGroup, i usecase.ItemsUseCase, l logger.Interface) {
	items := &itemsAPIRouter{
		items: i,
		l:     l,
	}

	h := handler.Group("/items")
	{
		h.GET("/:warehouse_id/quantity", items.getItemsQuantity)
		h.PUT("/", items.createItems)
	}
}

func (r *itemsAPIRouter) getItemsQuantity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("warehouse_id"))
	q, _ := r.items.Quantity(c.Request.Context(), id)

	c.JSON(http.StatusOK, q)
}

func (r *itemsAPIRouter) createItems(c *gin.Context) {
	var items []entity.Item
	if err := c.BindJSON(&items); err != nil {
		r.l.Error(err, "error binding JSON")
		errorResponse(c, http.StatusBadRequest, "provided data is invalid")
	}

	if err := r.items.CreateItems(c.Request.Context(), items); err != nil {
		r.l.Error(err, "failed to create item")
		errorResponse(c, http.StatusInternalServerError, "items service problems")

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

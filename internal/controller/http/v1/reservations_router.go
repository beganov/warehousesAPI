package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/robertgarayshin/warehousesAPI/internal/usecase"
	"github.com/robertgarayshin/warehousesAPI/pkg/custom_errors"
	"github.com/robertgarayshin/warehousesAPI/pkg/logger"
	"net/http"
)

type reservationsAPIRoutes struct {
	reservations usecase.ReservationsUsecase
	l            logger.Interface
}

func newReservationsAPIRoutes(handler *gin.RouterGroup, r usecase.ReservationsUsecase, l logger.Interface) {
	res := &reservationsAPIRoutes{
		reservations: r,
		l:            l,
	}

	h := handler.Group("/reserve")
	{
		h.POST("", res.reserve)
		h.DELETE("", res.deleteReservation)
	}
}

type reserveRequest struct {
	Ids []string `json:"ids"`
}

// @Summary     Reserve item
// @Description Reserve items in warehouse
// @ID          reserve
// @Tags  	    reservation
// @Accept      json
// @Produce     json
// @Param request body reserveRequest true "query params"
// @Success     201 {object} response
// @Failure     500 {object} response
// @Router      /reserve [post]
func (r *reservationsAPIRoutes) reserve(c *gin.Context) {
	var req reserveRequest
	if err := c.Bind(&req); err != nil {
		return
	}

	err := r.reservations.Reserve(c.Request.Context(), req.Ids)
	if errors.Is(err, custom_errors.ErrWarehouseUnavailable) {
		errorResponse(c, http.StatusForbidden, "warehouse is unavailable")

		return
	} else if err != nil {
		r.l.Error(err, "reservation error")
		errorResponse(c, http.StatusInternalServerError, "reservation error")

		return
	}

	successResponse(c, http.StatusCreated, "reservation successfully created")
}

// @Summary     Delete Item Reservation
// @Description Reserve items in warehouse
// @ID          deleteReservation
// @Tags  	    reservation
// @Accept      json
// @Produce     json
// @Param request body reserveRequest true "query params"
// @Success     200 {object} int
// @Failure		403 {object} response
// @Failure     500 {object} response
// @Router      /reserve [delete]
func (r *reservationsAPIRoutes) deleteReservation(c *gin.Context) {
	var request reserveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	err := r.reservations.CancelReservation(c.Request.Context(), request.Ids)
	if errors.Is(err, custom_errors.ErrNoReservation) {
		errorResponse(c, http.StatusForbidden, "item have no reservations")

		return
	} else if err != nil {
		r.l.Error(err, "error canceling reservation")
		errorResponse(c, http.StatusInternalServerError, "error canceling reservation")

		return
	}

	successResponse(c, http.StatusOK, "reservation successfully cancelled")
}

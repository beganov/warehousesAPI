package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehousesAPI/internal/usecase"
	"warehousesAPI/pkg/logger"
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

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /translation/history [get]
func (r *reservationsAPIRoutes) reserve(c *gin.Context) {
	var req reserveRequest
	if err := c.Bind(&req); err != nil {
		return
	}

	if err := r.reservations.Reserve(c.Request.Context(), req.Ids); err != nil {
		return
	}

	//c.JSON(http.StatusOK, historyResponse{translations})
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *reservationsAPIRoutes) deleteReservation(c *gin.Context) {
	var request reserveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	r.reservations.CancelReservation(c.Request.Context(), request.Ids)

	c.JSON(http.StatusOK, gin.H{})

	//
	//translation, err := r.t.Translate(
	//	c.Request.Context(),
	//	entity.Translation{
	//		Source:      request.Source,
	//		Destination: request.Destination,
	//		Original:    request.Original,
	//	},
	//)
	//if err != nil {
	//	r.l.Error(err, "http - v1 - doTranslate")
	//	errorResponse(c, http.StatusInternalServerError, "translation service problems")
	//
	//	return
	//}
	//
	//c.JSON(http.StatusOK, translation)
}

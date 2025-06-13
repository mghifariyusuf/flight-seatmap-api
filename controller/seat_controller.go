package controller

import (
	"flight-seatmap-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeatController struct {
	service service.SeatService
}

func NewSeatController(s service.SeatService) *SeatController {
	return &SeatController{service: s}
}

// GetAllSeats godoc
// @Summary Get all seats
// @Description Returns a list of all seat map data
// @Tags Seats
// @Produce json
// @Success 200 {array} model.Seat
// @Failure 500 {object} map[string]any
// @Router /seats [get]
func (c *SeatController) GetAllSeats(ctx *gin.Context) {
	seats, err := c.service.GetAllSeats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch seats"})
		return
	}
	ctx.JSON(http.StatusOK, seats)
}

type SeatSelectionRequest struct {
	Code string `json:"code" binding:"required"`
}

// SelectSeat godoc
// @Summary Select a seat
// @Description Select a seat by seat code
// @Tags Seats
// @Accept json
// @Produce json
// @Param input body SeatSelectionRequest true "Seat Code"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /seats/select [post]
func (c *SeatController) SelectSeat(ctx *gin.Context) {
	var req SeatSelectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := c.service.SelectSeat(req.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "seat selected"})
}

package routes

import (
	"flight-seatmap-api/controller"
	_ "flight-seatmap-api/docs"
	"flight-seatmap-api/repository"
	"flight-seatmap-api/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	seatRepo := repository.NewSeatRepository(db)
	seatSvc := service.NewSeatService(seatRepo)
	seatCtrl := controller.NewSeatController(seatSvc)

	api := r.Group("/seats")
	{
		api.GET("", seatCtrl.GetAllSeats)
		api.POST("/select", seatCtrl.SelectSeat)
	}

	// Swagger docs
	r.GET("/apidocs/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
}

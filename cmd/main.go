// @title Flight Seat Map API
// @version 1.0
// @description This is a backend API for selecting flight seats.
// @host localhost:8080
// @BasePath /
package main

import (
	"flight-seatmap-api/config"
	"flight-seatmap-api/db"
	"flight-seatmap-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	dbConn := config.ConnectDatabase(cfg)

	if err := db.SeedFromJSON("SeatMapResponse.json", dbConn); err != nil {
		logrus.WithError(err).Fatal("failed to seed database")
	}

	r := gin.Default()
	routes.RegisterRoutes(r, dbConn)
	r.Run(":8080")
}

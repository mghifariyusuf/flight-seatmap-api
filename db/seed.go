package db

import (
	"encoding/json"
	"flight-seatmap-api/model"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Structs for JSON parsing
type SeatMapResponse struct {
	SeatsItineraryParts []struct {
		SegmentSeatMaps []struct {
			PassengerSeatMaps []struct {
				SeatMap struct {
					Cabins []struct {
						Deck     string `json:"deck"`
						SeatRows []struct {
							RowNumber int `json:"rowNumber"`
							Seats     []struct {
								Code                   string   `json:"code"`
								Available              bool     `json:"available"`
								FreeOfCharge           bool     `json:"freeOfCharge"`
								RefundIndicator        string   `json:"refundIndicator"`
								RawSeatCharacteristics []string `json:"rawSeatCharacteristics"`
								Prices                 struct {
									Alternatives [][]struct {
										Amount   float64 `json:"amount"`
										Currency string  `json:"currency"`
									} `json:"alternatives"`
								} `json:"prices"`
							} `json:"seats"`
						} `json:"seatRows"`
					} `json:"cabins"`
				} `json:"seatMap"`
			} `json:"passengerSeatMaps"`
		} `json:"segmentSeatMaps"`
	} `json:"seatsItineraryParts"`
}

func SeedFromJSON(path string, db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Seat{}); err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var response SeatMapResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	for _, part := range response.SeatsItineraryParts {
		for _, segment := range part.SegmentSeatMaps {
			for _, pax := range segment.PassengerSeatMaps {
				for _, cabin := range pax.SeatMap.Cabins {
					for _, row := range cabin.SeatRows {
						for _, seat := range row.Seats {
							if seat.Code == "" {
								continue
							}

							price := 0.0
							currency := ""
							if len(seat.Prices.Alternatives) > 0 &&
								len(seat.Prices.Alternatives[0]) > 0 {
								price = seat.Prices.Alternatives[0][0].Amount
								currency = seat.Prices.Alternatives[0][0].Currency
							}

							s := model.Seat{
								Code:               seat.Code,
								RowNumber:          row.RowNumber,
								Cabin:              cabin.Deck,
								Available:          seat.Available,
								FreeOfCharge:       seat.FreeOfCharge,
								Price:              price,
								Currency:           currency,
								RefundIndicator:    seat.RefundIndicator,
								RawCharacteristics: strings.Join(seat.RawSeatCharacteristics, ","),
							}

							if err := db.Clauses(clause.OnConflict{
								UpdateAll: true,
							}).Create(&s).Error; err != nil {
								logrus.WithError(err).Errorf("Failed to insert or update seat %s.", s.Code)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

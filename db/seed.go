package db

import (
	"encoding/json"
	"flight-seatmap-api/model"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeatMapResponse struct {
	SeatsItineraryParts []SeatsItineraryPart `json:"seatsItineraryParts"`
}

type SeatsItineraryPart struct {
	SegmentSeatMaps []SegmentSeatMap `json:"segmentSeatMaps"`
}

type SegmentSeatMap struct {
	PassengerSeatMaps []PassengerSeatMap `json:"passengerSeatMaps"`
	Segment           any                `json:"segment"`
}

type PassengerSeatMap struct {
	SeatSelectionEnabledForPax bool    `json:"seatSelectionEnabledForPax"`
	SeatMap                    SeatMap `json:"seatMap"`
	Passenger                  any     `json:"passenger"`
}

type SeatMap struct {
	RowsDisabledCauses []string `json:"rowsDisabledCauses"`
	Aircraft           string   `json:"aircraft"`
	Cabins             []Cabin  `json:"cabins"`
}

type Cabin struct {
	Deck        string    `json:"deck"`
	SeatColumns []string  `json:"seatColumns"`
	SeatRows    []SeatRow `json:"seatRows"`
	FirstRow    int       `json:"firstRow"`
	LastRow     int       `json:"lastRow"`
}

type SeatRow struct {
	RowNumber int      `json:"rowNumber"`
	SeatCodes []string `json:"seatCodes"`
	Seats     []Seat   `json:"seats"`
}

type Seat struct {
	SlotCharacteristics    []string     `json:"slotCharacteristics"`
	StorefrontSlotCode     string       `json:"storefrontSlotCode"`
	Available              bool         `json:"available"`
	Code                   string       `json:"code"`
	Designations           []string     `json:"designations"`
	Entitled               bool         `json:"entitled"`
	FeeWaived              bool         `json:"feeWaived"`
	EntitledRuleID         string       `json:"entitledRuleId"`
	FeeWaivedRuleID        string       `json:"feeWaivedRuleId"`
	SeatCharacteristics    []string     `json:"seatCharacteristics"`
	Limitations            []string     `json:"limitations"`
	RefundIndicator        string       `json:"refundIndicator"`
	FreeOfCharge           bool         `json:"freeOfCharge"`
	Prices                 SeatPriceSet `json:"prices"`
	Taxes                  SeatPriceSet `json:"taxes"`
	Total                  SeatPriceSet `json:"total"`
	OriginallySelected     bool         `json:"originallySelected"`
	RawSeatCharacteristics []string     `json:"rawSeatCharacteristics"`
}

type SeatPriceSet struct {
	Alternatives [][]SeatPrice `json:"alternatives"`
}

type SeatPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
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
							prices := ""
							for _, alt := range seat.Prices.Alternatives {
								for _, price := range alt {
									prices += price.Currency + ":" + fmt.Sprintf("%.2f", price.Amount) + ","
								}
							}

							taxes := ""
							for _, alt := range seat.Taxes.Alternatives {
								for _, taxe := range alt {
									taxes += taxe.Currency + ":" + fmt.Sprintf("%.2f", taxe.Amount) + ","
								}
							}

							total := ""
							for _, alt := range seat.Total.Alternatives {
								for _, t := range alt {
									total += t.Currency + ":" + fmt.Sprintf("%.2f", t.Amount) + ","
								}
							}

							seat := model.Seat{
								RowNumber:              row.RowNumber,
								SlotCharacteristics:    strings.Join(seat.SlotCharacteristics, ","),
								StorefrontSlotCode:     seat.StorefrontSlotCode,
								Available:              seat.Available,
								Code:                   seat.Code,
								Designations:           strings.Join(seat.Designations, ","),
								Entitled:               seat.Entitled,
								FeeWaived:              seat.FeeWaived,
								EntitledRuleID:         seat.EntitledRuleID,
								FeeWaivedRuleID:        seat.FeeWaivedRuleID,
								SeatCharacteristics:    strings.Join(seat.SeatCharacteristics, ","),
								Limitations:            strings.Join(seat.Limitations, ","),
								RefundIndicator:        seat.RefundIndicator,
								FreeOfCharge:           seat.FreeOfCharge,
								Prices:                 prices,
								Taxes:                  taxes,
								Total:                  total,
								OriginallySelected:     seat.OriginallySelected,
								RawSeatCharacteristics: strings.Join(seat.RawSeatCharacteristics, ","),
							}

							err := db.Clauses(clause.OnConflict{
								Columns: []clause.Column{
									{Name: "row_number"},
									{Name: "code"},
									{Name: "slot_characteristics"},
								},
								UpdateAll: true,
							}).Create(&seat).Error
							if err != nil {
								logrus.WithError(err).Errorf("Failed to insert")
							}
						}
					}
				}
			}
		}
	}

	return nil
}

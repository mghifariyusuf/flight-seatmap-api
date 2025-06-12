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

type SeatMapResponse struct {
	SeatsItineraryParts []SeatsItineraryPart `json:"seatsItineraryParts"`
}

type SeatsItineraryPart struct {
	SegmentSeatMaps []SegmentSeatMap `json:"segmentSeatMaps"`
}

type SegmentSeatMap struct {
	PassengerSeatMaps []PassengerSeatMap `json:"passengerSeatMaps"`
}

type PassengerSeatMap struct {
	SeatMap SeatMap `json:"seatMap"`
}

type SeatMap struct {
	Cabins []Cabin `json:"cabins"`
}

type Cabin struct {
	Deck     string    `json:"deck"`
	SeatRows []SeatRow `json:"seatRows"`
}

type SeatRow struct {
	RowNumber int    `json:"rowNumber"`
	Seats     []Seat `json:"seats"`
	Slots     []Slot `json:"slots"`
}

type Seat struct {
	Code                   string       `json:"code"`
	Available              bool         `json:"available"`
	FreeOfCharge           bool         `json:"freeOfCharge"`
	RefundIndicator        string       `json:"refundIndicator"`
	RawSeatCharacteristics []string     `json:"rawSeatCharacteristics"`
	Prices                 SeatPriceSet `json:"prices"`
}

type SeatPriceSet struct {
	Alternatives [][]SeatPrice `json:"alternatives"`
}

type SeatPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Slot struct {
	SlotCharacteristics []string `json:"slotCharacteristics"`
	StorefrontSlotCode  string   `json:"storefrontSlotCode"`
	Available           bool     `json:"available"`
	Entitled            bool     `json:"entitled"`
	FeeWaived           bool     `json:"feeWaived"`
	FreeOfCharge        bool     `json:"freeOfCharge"`
	OriginallySelected  bool     `json:"originallySelected"`
}

func SeedFromJSON(path string, db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Slot{}); err != nil {
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
							code := seat.Code
							slotType := "seat"
							if code == "" {
								code = "BLANK"
								slotType = "blank"
							}

							price := 0.0
							currency := ""
							if len(seat.Prices.Alternatives) > 0 && len(seat.Prices.Alternatives[0]) > 0 {
								price = seat.Prices.Alternatives[0][0].Amount
								currency = seat.Prices.Alternatives[0][0].Currency
							}

							slot := model.Slot{
								Code:                code,
								RowNumber:           row.RowNumber,
								Cabin:               cabin.Deck,
								Available:           seat.Available,
								FreeOfCharge:        seat.FreeOfCharge,
								FeeWaived:           false,
								Entitled:            false,
								OriginallySelected:  false,
								Type:                slotType,
								SlotCharacteristics: "",
								Currency:            currency,
								Price:               price,
								RefundIndicator:     seat.RefundIndicator,
								RawCharacteristics:  strings.Join(seat.RawSeatCharacteristics, ","),
							}

							err := db.Clauses(clause.OnConflict{
								Columns: []clause.Column{{Name: "code"}, {Name: "row_number"}},
								DoUpdates: clause.AssignmentColumns([]string{
									"cabin", "available", "free_of_charge", "currency", "price", "refund_indicator",
									"raw_characteristics", "type",
								}),
							}).Create(&slot).Error
							if err != nil {
								logrus.WithError(err).Errorf("Failed to insert or update slot %s.", code)
							}
						}

						for _, slot := range row.Slots {
							code := slot.StorefrontSlotCode
							slotType := "blank"
							if code == "" {
								code = "BLANK"
							}

							newSlot := model.Slot{
								Code:                code,
								RowNumber:           row.RowNumber,
								Cabin:               cabin.Deck,
								Available:           slot.Available,
								FreeOfCharge:        slot.FreeOfCharge,
								FeeWaived:           slot.FeeWaived,
								Entitled:            slot.Entitled,
								OriginallySelected:  slot.OriginallySelected,
								Type:                slotType,
								SlotCharacteristics: strings.Join(slot.SlotCharacteristics, ","),
							}

							err := db.Clauses(clause.OnConflict{
								Columns: []clause.Column{{Name: "code"}, {Name: "row_number"}},
								DoUpdates: clause.AssignmentColumns([]string{
									"available", "free_of_charge", "fee_waived", "entitled", "originally_selected",
									"type", "slot_characteristics",
								}),
							}).Create(&newSlot).Error
							if err != nil {
								logrus.WithError(err).Errorf("Failed to insert or update slot %s.", code)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

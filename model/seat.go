package model

type Seat struct {
	RowNumber              int    `gorm:"primaryKey"`
	SlotCharacteristics    string `gorm:"primaryKey"`
	StorefrontSlotCode     string
	Available              bool
	Code                   string `gorm:"primaryKey"`
	Designations           string
	Entitled               bool
	FeeWaived              bool
	EntitledRuleID         string
	FeeWaivedRuleID        string
	SeatCharacteristics    string
	Limitations            string
	RefundIndicator        string
	FreeOfCharge           bool
	Prices                 string
	Taxes                  string
	Total                  string
	OriginallySelected     bool
	RawSeatCharacteristics string
}

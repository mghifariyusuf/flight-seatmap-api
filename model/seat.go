package model

type Seat struct {
	RowNumber              int
	SlotCharacteristics    string
	StorefrontSlotCode     string
	Available              bool
	Code                   string
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

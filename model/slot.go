package model

type Slot struct {
	Code                string `gorm:"index"`
	RowNumber           int
	Cabin               string
	Available           bool
	FreeOfCharge        bool
	FeeWaived           bool
	Entitled            bool
	OriginallySelected  bool
	Type                string
	SlotCharacteristics string
	Currency            string
	Price               float64
	RefundIndicator     string
	RawCharacteristics  string
}

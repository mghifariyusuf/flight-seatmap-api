package model

// Seat defines the seat data structure stored in PostgreSQL
type Seat struct {
	Code               string `gorm:"primarykey"`
	RowNumber          int
	Cabin              string
	Available          bool
	FreeOfCharge       bool
	Currency           string
	Price              float64
	RefundIndicator    string
	RawCharacteristics string
}

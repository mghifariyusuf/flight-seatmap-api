package repository

import (
	"errors"
	"flight-seatmap-api/model"

	"gorm.io/gorm"
)

type SeatRepository interface {
	GetAll() ([]model.Slot, error)
	SelectSeat(code string) error
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db}
}

func (r *seatRepository) GetAll() ([]model.Slot, error) {
	var seats []model.Slot
	err := r.db.Find(&seats).Error
	return seats, err
}

func (r *seatRepository) SelectSeat(code string) error {
	res := r.db.Model(&model.Slot{}).
		Where("code = ? AND available = true", code).
		Update("available", false)
	if res.RowsAffected == 0 {
		return errors.New("seat not available or does not exist")
	}
	return res.Error
}

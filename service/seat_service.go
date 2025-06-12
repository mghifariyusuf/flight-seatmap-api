package service

import (
	"errors"
	"flight-seatmap-api/model"
	"flight-seatmap-api/repository"
)

type SeatService interface {
	GetAllSeats() ([]model.Seat, error)
	SelectSeat(code string) error
}

type seatService struct {
	repo repository.SeatRepository
}

func NewSeatService(repo repository.SeatRepository) SeatService {
	return &seatService{repo}
}

func (s *seatService) GetAllSeats() ([]model.Seat, error) {
	return s.repo.GetAll()
}

func (s *seatService) SelectSeat(code string) error {
	if code == "" {
		return errors.New("seat code cannot be empty")
	}
	return s.repo.SelectSeat(code)
}

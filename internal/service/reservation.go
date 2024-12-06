package service

import (
	"context"
	"user_balance/internal/entity"
	"user_balance/internal/repository"
)

type ReservationService struct {
	repo repository.Reservation
}

func NewReservationService(repo repository.Reservation) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	id, err := s.repo.CreateReservation(ctx, reservation)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ReservationService) GetReservation(ctx context.Context, reservationID int) (entity.Reservation, error) {
	return s.repo.GetReservation(ctx, reservationID)
}

func (s *ReservationService) RefundReservation(ctx context.Context, reservationId int) error {
	return s.repo.RefundReservation(ctx, reservationId)
}

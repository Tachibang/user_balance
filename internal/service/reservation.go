package service

import (
	"context"
	"user_balance/internal/entity"
	"user_balance/internal/repository"

	"github.com/sirupsen/logrus"
)

type ReservationService struct {
	repo   repository.Reservation
	logger *logrus.Logger
}

func NewReservationService(repo repository.Reservation, logger *logrus.Logger) *ReservationService {
	return &ReservationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ReservationService) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	s.logger.Infof("Создание резервации: %+v", reservation)
	id, err := s.repo.CreateReservation(ctx, reservation)
	if err != nil {
		s.logger.Errorf("Ошибка при создании резервации: %w", err)
		return 0, err
	}
	s.logger.Infof("Резервация создана с ID: %d", id)
	return id, nil
}

func (s *ReservationService) GetReservation(ctx context.Context, reservationID int) (entity.Reservation, error) {
	s.logger.Infof("Получение резервации с ID: %d", reservationID)
	reservation, err := s.repo.GetReservation(ctx, reservationID)
	if err != nil {
		s.logger.Errorf("Ошибка при получении резервации с ID %d: %w", reservationID, err)
		return entity.Reservation{}, err
	}
	s.logger.Infof("Резервация с ID %d успешно получена: %+v", reservationID, reservation)
	return reservation, nil
}

func (s *ReservationService) RefundReservation(ctx context.Context, reservationId int) error {
	s.logger.Infof("Возврат резервации с ID: %d", reservationId)
	if err := s.repo.RefundReservation(ctx, reservationId); err != nil {
		s.logger.Errorf("Ошибка при возврате резервации с ID %d: %w", reservationId, err)
		return err
	}
	s.logger.Infof("Резервация с ID %d успешно возвращена", reservationId)
	return nil
}

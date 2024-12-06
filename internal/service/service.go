package service

import (
	"context"
	"user_balance/internal/entity"
	"user_balance/internal/repository"

	"github.com/sirupsen/logrus"
)

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
	GetAccount(ctx context.Context, id int) (entity.Account, error)
	Deposit(ctx context.Context, id, amount int) (int, int, error)
	Withdraw(ctx context.Context, id, amount int) (int, int, error)
	Transfer(ctx context.Context, fromID, toID, amount int) (int, int, error)
}

type Reservation interface {
	CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error)
	GetReservation(ctx context.Context, reservationID int) (entity.Reservation, error)
	RefundReservation(ctx context.Context, reservationId int) error
}

type Product interface {
	CreateProduct(ctx context.Context, name string) (int, error)
	GetProduct(ctx context.Context, id int) (entity.Product, error)
}

type Service struct {
	Account
	Reservation
	Product
}

func NewService(repository *repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		Account:     NewAccountService(repository),
		Reservation: NewReservationService(repository),
		Product:     NewProductService(repository, log),
	}
}

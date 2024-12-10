package repository

import (
	"context"
	"database/sql"
	"time"
	"user_balance/internal/entity"
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

type Operation interface {
	GetMonthlyOperations(ctx context.Context, startDate, endDate time.Time) ([]entity.Operation, error)
}

type Repository struct {
	Account
	Product
	Reservation
	Operation
}

func NewRepository(pg *sql.DB) *Repository {
	return &Repository{
		Account:     NewAccountRepo(pg),
		Product:     NewProductRepo(pg),
		Reservation: NewReservationRepo(pg),
		Operation:   NewOperationRepo(pg),
	}
}

package repository

import (
	"context"
	"database/sql"
	"time"
	"user_balance/internal/entity"
	"user_balance/internal/repository/repoerrs"
)

type ReservationRepo struct {
	pg *sql.DB
}

func NewReservationRepo(pg *sql.DB) *ReservationRepo {
	return &ReservationRepo{pg}
}

func (r *ReservationRepo) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var balance int
	var accountDeletedAt *time.Time
	queryCheckAccount := "SELECT balance, deleted_at FROM accounts WHERE id = $1"
	err = tx.QueryRowContext(ctx, queryCheckAccount, reservation.AccountId).Scan(&balance, &accountDeletedAt)

	if err != nil {
		return 0, err
	}
	if accountDeletedAt != nil {
		return 0, repoerrs.ErrDataDeleted
	}
	if balance < reservation.Amount {
		return 0, repoerrs.ErrNotEnoughBalance
	}

	var productDeletedAt *time.Time
	queryCheckProduct := "SELECT deleted_at FROM products WHERE id = $1"
	err = tx.QueryRowContext(ctx, queryCheckProduct, reservation.ProductId).Scan(&productDeletedAt)
	if err != nil {
		return 0, err
	}
	if productDeletedAt != nil {
		return 0, repoerrs.ErrDataDeleted
	}

	queryUpdateBalance := `
		UPDATE accounts
		SET balance = balance - $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, reservation.Amount, reservation.AccountId)
	if err != nil {
		return 0, err
	}

	var reservationID int
	queryInsertReservation := `
		INSERT INTO reservations (account_id, product_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, queryInsertReservation,
		reservation.AccountId, reservation.ProductId, reservation.Amount,
	).Scan(&reservationID)

	if err != nil {
		return 0, err
	}

	queryInsertOperation := `
		INSERT INTO operations (account_id, amount, operation_type, product_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(ctx, queryInsertOperation,
		reservation.AccountId, reservation.Amount, "reservation", reservation.ProductId,
	)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return reservationID, nil
}

func (r *ReservationRepo) GetReservation(ctx context.Context, reservationID int) (entity.Reservation, error) {
	queryGetReservation := `
    SELECT id, account_id, product_id, amount, created_at, deleted_at
    FROM reservations
    WHERE id = $1
    `

	var reservation entity.Reservation
	err := r.pg.QueryRowContext(ctx, queryGetReservation, reservationID).Scan(
		&reservation.Id,
		&reservation.AccountId,
		&reservation.ProductId,
		&reservation.Amount,
		&reservation.CreatedAt,
		&reservation.DeletedAt,
	)

	if err != nil {
		return reservation, err
	}

	if reservation.DeletedAt != nil {
		return reservation, repoerrs.ErrDataDeleted
	}

	return reservation, nil
}

func (r *ReservationRepo) RefundReservation(ctx context.Context, reservationId int) error {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queryGetReservation := `
	SELECT account_id, amount, product_id, created_at, deleted_at FROM reservations WHERE id = $1
	`
	var reservation entity.Reservation
	err = tx.QueryRowContext(ctx, queryGetReservation, reservationId).Scan(
		&reservation.AccountId,
		&reservation.Amount,
		&reservation.ProductId,
		&reservation.CreatedAt,
		&reservation.DeletedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if reservation.DeletedAt != nil {
		tx.Rollback()
		return err
	}

	queryGetAccount := `
	SELECT deleted_at FROM accounts WHERE id = $1
	`
	var accountDeletedAt *time.Time
	err = tx.QueryRowContext(ctx, queryGetAccount, reservation.AccountId).Scan(&accountDeletedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	if accountDeletedAt != nil {
		tx.Rollback()
		return err
	}

	queryUpdateReservation := `
	UPDATE reservations
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, queryUpdateReservation, reservationId)
	if err != nil {
		tx.Rollback()
		return err
	}

	queryUpdateBalance := `
	UPDATE accounts
	SET balance = balance + $1, updated_at = NOW()
	WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, reservation.Amount, reservation.AccountId)
	if err != nil {
		tx.Rollback()
		return err
	}

	queryInsertOperation := `
	INSERT INTO operations (account_id, amount, operation_type, product_id)
	VALUES ($1, $2, $3, $4)
	`
	_, err = tx.ExecContext(ctx, queryInsertOperation,
		reservation.AccountId, reservation.Amount, "refund", reservation.ProductId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

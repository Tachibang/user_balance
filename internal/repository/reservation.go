package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"user_balance/internal/entity"
)

type ReservationRepo struct {
	pg *sql.DB
}

func NewReservationRepo(pg *sql.DB) *ReservationRepo {
	return &ReservationRepo{pg}
}

func (r *ReservationRepo) CreateReservation(ctx context.Context, reservation entity.Reservation) (int, error) {
	log.Printf("Создание резервации для аккаунта ID: %d, сумма: %d", reservation.AccountId, reservation.Amount)
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v", err)
		return 0, errors.New("не удалось начать транзакцию")
	}

	queryGetBalance := `
    SELECT balance FROM accounts WHERE id = $1
    `
	var balance int
	err = tx.QueryRowContext(ctx, queryGetBalance, reservation.AccountId).Scan(&balance)
	if err != nil {
		log.Printf("Ошибка при получении баланса аккаунта ID %d: %v", reservation.AccountId, err)
		return 0, errors.New("не удалось получить баланс аккаунта")
	}

	if balance < reservation.Amount {
		log.Printf("Недостаточно средств на счете для резервации. Баланс: %d, требуется: %d", balance, reservation.Amount)
		return 0, errors.New("недостаточно средств")
	}

	queryUpdateBalance := `
    UPDATE accounts
    SET balance = balance - $1
    WHERE id = $2
    `
	_, err = tx.ExecContext(ctx, queryUpdateBalance, reservation.Amount, reservation.AccountId)
	if err != nil {
		log.Printf("Ошибка при обновлении баланса аккаунта ID %d: %v", reservation.AccountId, err)
		return 0, errors.New("не удалось обновить баланс аккаунта")
	}

	queryInsertReservation := `
    INSERT INTO reservations (account_id, product_id, amount)
    VALUES ($1, $2, $3)
    RETURNING id
    `
	var reservationID int
	err = tx.QueryRowContext(ctx, queryInsertReservation,
		reservation.AccountId, reservation.ProductId, reservation.Amount,
	).Scan(&reservationID)
	if err != nil {
		log.Printf("Ошибка при создании резервации: %v", err)
		return 0, errors.New("не удалось создать резервацию")
	}

	queryInsertOperation := `
    INSERT INTO operations (account_id, amount, operation_type, product_id, created_at)
    VALUES ($1, $2, $3, $4, NOW())
    `
	_, err = tx.ExecContext(ctx, queryInsertOperation,
		reservation.AccountId, reservation.Amount, "reservation", reservation.ProductId,
	)
	if err != nil {
		log.Printf("Ошибка при добавлении операции: %v", err)
		return 0, errors.New("не удалось добавить операцию")
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return 0, errors.New("не удалось зафиксировать транзакцию")
	}

	log.Printf("Резервация успешно создана с ID: %d", reservationID)
	return reservationID, nil
}

func (r *ReservationRepo) GetReservation(ctx context.Context, reservationID int) (entity.Reservation, error) {
	log.Printf("Получение резервации с ID: %d", reservationID)
	queryGetReservation := `
    SELECT *
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Резервация с ID %d не найдена", reservationID)
			return reservation, errors.New("резервация не найдена")
		}
		log.Printf("Ошибка при получении резервации: %v", err)
		return reservation, errors.New("не удалось получить резервацию")
	}

	log.Printf("Резервация получена: ID %d, аккаунт ID %d, сумма: %d", reservation.Id, reservation.AccountId, reservation.Amount)
	return reservation, nil
}

func (r *ReservationRepo) RefundReservation(ctx context.Context, reservationId int) error {
	log.Printf("Возврат резервации с ID: %d", reservationId)
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v", err)
		return errors.New("не удалось начать транзакцию")
	}

	queryGetReservation := `
    SELECT account_id, amount, created_at FROM reservations WHERE id = $1
    `
	var reservation entity.Reservation
	err = tx.QueryRowContext(ctx, queryGetReservation, reservationId).Scan(&reservation.AccountId, &reservation.Amount, &reservation.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Резервация с ID %d не найдена", reservationId)
			return errors.New("резервация не найдена")
		}
		log.Printf("Ошибка при получении резервации: %v", err)
		return errors.New("не удалось получить резервацию")
	}

	queryUpdateReservation := `
    DELETE FROM reservations WHERE id = $1
    `
	_, err = tx.ExecContext(ctx, queryUpdateReservation, reservationId)
	if err != nil {
		log.Printf("Ошибка при удалении резервации: %v", err)
		return errors.New("не удалось удалить резервацию")
	}

	queryUpdateBalance := `
    UPDATE accounts
    SET balance = balance + $1
    WHERE id = $2
    `
	_, err = tx.ExecContext(ctx, queryUpdateBalance, reservation.Amount, reservation.AccountId)
	if err != nil {
		log.Printf("Ошибка при обновлении баланса аккаунта ID %d: %v", reservation.AccountId, err)
		return errors.New("не удалось обновить баланс аккаунта")
	}

	queryInsertOperation := `
    INSERT INTO operations (account_id, amount, operation_type, product_id, created_at)
    VALUES ($1, $2, $3, $4, NOW())
    `
	_, err = tx.ExecContext(ctx, queryInsertOperation,
		reservation.AccountId, reservation.Amount, "refund", reservation.ProductId,
	)
	if err != nil {
		log.Printf("Ошибка при добавлении операции: %v", err)
		return errors.New("не удалось добавить операцию")
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка при коммите транзакции: %v", err)
		return errors.New("не удалось зафиксировать транзакцию")
	}

	log.Printf("Возврат резервации с ID %d успешно выполнен", reservationId)
	return nil
}

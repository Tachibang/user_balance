package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"user_balance/internal/entity"
)

type AccountRepo struct {
	pg *sql.DB
}

func NewAccountRepo(pg *sql.DB) *AccountRepo {
	return &AccountRepo{pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context) (int, error) {
	var id int
	query := "INSERT INTO accounts (balance) VALUES ($1) RETURNING id"
	err := r.pg.QueryRowContext(ctx, query, 0).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при создании счета: %v\n", err)
		return 0, errors.New("не удалось создать счет")
	}
	log.Printf("Счет с ID %d успешно создан\n", id)
	return id, nil
}

func (r *AccountRepo) GetAccount(ctx context.Context, id int) (entity.Account, error) {
	var account entity.Account
	query := "SELECT * FROM accounts WHERE id = $1"
	err := r.pg.QueryRowContext(ctx, query, id).Scan(
		&account.Id,
		&account.Balance,
		&account.CreatedAt)
	if err != nil {
		log.Printf("Ошибка при получении счета с ID %d: %v\n", id, err)
		return entity.Account{}, errors.New("не удалось получить счет")
	}
	log.Printf("Счет с ID %d успешно получен\n", id)
	return account, nil
}

func (r *AccountRepo) Deposit(ctx context.Context, id, amount int) (int, int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v\n", err)
		return 0, 0, errors.New("не удалось начать транзакцию")
	}

	queryUpdateBalance := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance
	`

	var newBalance int
	err = tx.QueryRowContext(ctx, queryUpdateBalance, amount, id).Scan(&newBalance)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при обновлении баланса счета: %v\n", err)
		return 0, 0, errors.New("не удалось обновить баланс")
	}

	queryInsertOperation := `
		INSERT INTO operations (account_id, amount, operation_type, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = tx.ExecContext(ctx, queryInsertOperation, id, amount, "deposit")
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи операции: %v\n", err)
		return 0, 0, errors.New("не удалось записать операцию")
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка при коммите транзакции: %v\n", err)
		return 0, 0, errors.New("не удалось завершить транзакцию")
	}

	log.Printf("Внесено %d на счет с ID %d, новый баланс: %d\n", amount, id, newBalance)
	return id, newBalance, nil
}

func (r *AccountRepo) Withdraw(ctx context.Context, id, amount int) (int, int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v\n", err)
		return 0, 0, errors.New("не удалось начать транзакцию")
	}

	queryGetBalance := `
	SELECT balance FROM accounts WHERE id=$1
	`

	var balance int
	err = tx.QueryRowContext(ctx, queryGetBalance, id).Scan(&balance)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при получении баланса счета: %v\n", err)
		return 0, 0, errors.New("не удалось получить баланс")
	}

	if balance < amount {
		tx.Rollback()
		log.Println("Ошибка: недостаточно средств")
		return 0, 0, errors.New("недостаточно средств")
	}

	queryUpdateBalance := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		RETURNING balance
	`

	var newBalance int
	err = tx.QueryRowContext(ctx, queryUpdateBalance, amount, id).Scan(&newBalance)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при обновлении баланса счета: %v\n", err)
		return 0, 0, errors.New("не удалось обновить баланс")
	}

	queryInsertOperation := `
		INSERT INTO operations (account_id, amount, operation_type, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = tx.ExecContext(ctx, queryInsertOperation, id, amount, "withdraw")
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи операции: %v\n", err)
		return 0, 0, errors.New("не удалось записать операцию")
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка при коммите транзакции: %v\n", err)
		return 0, 0, errors.New("не удалось завершить транзакцию")
	}

	log.Printf("Снято %d с счета с ID %d, новый баланс: %d\n", amount, id, newBalance)
	return id, newBalance, nil
}

func (r *AccountRepo) Transfer(ctx context.Context, fromID, toID, amount int) (int, int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка при начале транзакции: %v\n", err)
		return 0, 0, errors.New("не удалось начать транзакцию")
	}

	queryGetBalance := `
    SELECT balance FROM accounts WHERE id=$1
    `

	var fromBalance int
	err = tx.QueryRowContext(ctx, queryGetBalance, fromID).Scan(&fromBalance)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при получении баланса с счета ID %d: %v\n", fromID, err)
		return 0, 0, errors.New("не удалось получить баланс")
	}

	if fromBalance < amount {
		tx.Rollback()
		log.Println("Ошибка: недостаточно средств")
		return 0, 0, errors.New("недостаточно средств")
	}

	queryUpdateFromBalance := `
    UPDATE accounts
    SET balance = balance - $1
    WHERE id = $2
    RETURNING balance
    `

	var newFromBalance int
	err = tx.QueryRowContext(ctx, queryUpdateFromBalance, amount, fromID).Scan(&newFromBalance)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при обновлении баланса счета ID %d: %v\n", fromID, err)
		return 0, 0, errors.New("не удалось обновить баланс источника")
	}

	queryUpdateToBalance := `
    UPDATE accounts
    SET balance = balance + $1
    WHERE id = $2
    RETURNING balance
    `

	var newToBalance int
	err = tx.QueryRowContext(ctx, queryUpdateToBalance, amount, toID).Scan(&newToBalance)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при обновлении баланса счета ID %d: %v\n", toID, err)
		return 0, 0, errors.New("не удалось обновить баланс получателя")
	}

	queryInsertFromOperation := `
    INSERT INTO operations (account_id, amount, operation_type, created_at)
    VALUES ($1, $2, $3, NOW())
    `
	_, err = tx.ExecContext(ctx, queryInsertFromOperation, fromID, amount, "transfer_out")
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи операции перевода из счета ID %d: %v\n", fromID, err)
		return 0, 0, errors.New("не удалось записать операцию перевода из")
	}

	queryInsertToOperation := `
    INSERT INTO operations (account_id, amount, operation_type, created_at)
    VALUES ($1, $2, $3, NOW())
    `
	_, err = tx.ExecContext(ctx, queryInsertToOperation, toID, amount, "transfer_in")
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи операции перевода на счет ID %d: %v\n", toID, err)
		return 0, 0, errors.New("не удалось записать операцию перевода на")
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка при коммите транзакции: %v\n", err)
		return 0, 0, errors.New("не удалось завершить транзакцию")
	}

	log.Printf("Переведено %d с счета ID %d на счет ID %d, новые балансы: %d и %d\n", amount, fromID, toID, newFromBalance, newToBalance)
	return newFromBalance, newToBalance, nil
}

package repository

import (
	"context"
	"database/sql"
	"time"
	"user_balance/internal/entity"
	"user_balance/internal/repository/repoerrs"
)

type AccountRepo struct {
	pg *sql.DB
}

func NewAccountRepo(pg *sql.DB) *AccountRepo {
	return &AccountRepo{pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context) (int, error) {
	var id int
	query := "INSERT INTO accounts DEFAULT VALUES RETURNING id"
	err := r.pg.QueryRowContext(ctx, query).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountRepo) GetAccount(ctx context.Context, id int) (entity.Account, error) {
	var account entity.Account
	query := `
		SELECT id, balance, created_at, updated_at, deleted_at
		FROM accounts
		WHERE id = $1
	`

	err := r.pg.QueryRowContext(ctx, query, id).Scan(
		&account.Id,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Account{}, sql.ErrNoRows
		}
		return entity.Account{}, err
	}

	if account.DeletedAt != nil {
		return entity.Account{}, repoerrs.ErrDataDeleted
	}

	return account, nil
}

func (r *AccountRepo) Deposit(ctx context.Context, id, amount int) (int, int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	queryUpdateBalance := `
		UPDATE accounts
		SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2
		RETURNING balance, deleted_at
	`

	var newBalance int
	var deletedCheck *time.Time
	err = tx.QueryRowContext(ctx, queryUpdateBalance, amount, id).Scan(&newBalance, &deletedCheck)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	if deletedCheck != nil {
		return 0, 0, repoerrs.ErrDataDeleted
	}

	queryInsertOperation := `
		INSERT INTO operations (account_id, amount, operation_type, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = tx.ExecContext(ctx, queryInsertOperation, id, amount, "deposit")
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	return id, newBalance, nil
}

func (r *AccountRepo) Withdraw(ctx context.Context, id, amount int) (int, int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	queryGetBalance := `
	SELECT balance, deleted_at FROM accounts WHERE id=$1
	`

	var balance int
	var deletedAt *time.Time
	err = tx.QueryRowContext(ctx, queryGetBalance, id).Scan(&balance, &deletedAt)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	if deletedAt != nil {
		tx.Rollback()
		return 0, 0, repoerrs.ErrNotEnoughBalance
	}

	if balance < amount {
		tx.Rollback()
		return 0, 0, repoerrs.ErrNotEnoughBalance
	}

	queryUpdateBalance := `
		UPDATE accounts
		SET balance = balance - $1, updated_at = now()
		WHERE id = $2
		RETURNING balance
	`

	var newBalance int
	err = tx.QueryRowContext(ctx, queryUpdateBalance, amount, id).Scan(&newBalance)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	queryInsertOperation := `
		INSERT INTO operations (account_id, amount, operation_type, created_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err = tx.ExecContext(ctx, queryInsertOperation, id, amount, "withdraw")
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	return id, newBalance, nil
}

func (r *AccountRepo) Transfer(ctx context.Context, fromID, toID, amount int) (int, int, error) {
	tx, err := r.pg.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	queryGetBalance := `
    SELECT balance, deleted_at FROM accounts WHERE id=$1
    `

	var fromBalance int
	var fromDeletedAt *time.Time
	err = tx.QueryRowContext(ctx, queryGetBalance, fromID).Scan(&fromBalance, &fromDeletedAt)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}
	if fromDeletedAt != nil {
		tx.Rollback()
		return 0, 0, repoerrs.ErrNotEnoughBalance
	}

	var toBalance int
	var toDeletedAt *time.Time
	err = tx.QueryRowContext(ctx, queryGetBalance, toID).Scan(&toBalance, &toDeletedAt)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	if toDeletedAt != nil {
		tx.Rollback()
		return 0, 0, repoerrs.ErrNotEnoughBalance
	}

	if fromBalance < amount {
		tx.Rollback()
		return 0, 0, repoerrs.ErrNotEnoughBalance
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
		return 0, 0, err
	}

	queryUpdateToBalance := `
    UPDATE accounts
    SET balance = balance + $1, updated_at = now()
    WHERE id = $2
    RETURNING balance
    `

	var newToBalance int
	err = tx.QueryRowContext(ctx, queryUpdateToBalance, amount, toID).Scan(&newToBalance)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	queryInsertFromOperation := `
    INSERT INTO operations (account_id, amount, operation_type, created_at)
    VALUES ($1, $2, $3, NOW())
    `

	_, err = tx.ExecContext(ctx, queryInsertFromOperation, fromID, amount, "transfer_out")
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	queryInsertToOperation := `
    INSERT INTO operations (account_id, amount, operation_type, created_at)
    VALUES ($1, $2, $3, NOW())
    `

	_, err = tx.ExecContext(ctx, queryInsertToOperation, toID, amount, "transfer_in")
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	return newFromBalance, newToBalance, nil
}

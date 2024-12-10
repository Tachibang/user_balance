package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"user_balance/internal/entity"
)

type OperationRepo struct {
	pg *sql.DB
}

func NewOperationRepo(pg *sql.DB) *OperationRepo {
	return &OperationRepo{pg}
}

func (r *OperationRepo) GetMonthlyOperations(ctx context.Context, startDate, endDate time.Time) ([]entity.Operation, error) {
	startDate = startDate.UTC()
	endDate = endDate.UTC()

	query := `
		SELECT 
			id, account_id, amount, operation_type, product_id, description, created_at, updated_at, deleted_at
		FROM operations 
		WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL
	`

	log.Printf("Start Date: %s, End Date: %s", startDate, endDate)

	rows, err := r.pg.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []entity.Operation
	for rows.Next() {
		var op entity.Operation
		err := rows.Scan(
			&op.Id,
			&op.AccountId,
			&op.Amount,
			&op.OperationType,
			&op.ProductId,
			&op.Description,
			&op.CreatedAt,
			&op.UpdatedAt,
			&op.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, op)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return operations, nil
}

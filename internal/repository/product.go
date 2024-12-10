package repository

import (
	"context"
	"database/sql"
	"user_balance/internal/entity"
	"user_balance/internal/repository/repoerrs"
)

type ProductRepo struct {
	pg *sql.DB
}

func NewProductRepo(pg *sql.DB) *ProductRepo {
	return &ProductRepo{pg}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, name string) (int, error) {
	var id int
	query := "INSERT INTO products (name) VALUES ($1) RETURNING id"
	err := r.pg.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductRepo) GetProduct(ctx context.Context, id int) (entity.Product, error) {
	var product entity.Product
	query := "SELECT id, name, deleted_at FROM products WHERE id = $1"
	err := r.pg.QueryRowContext(ctx, query, id).Scan(
		&product.Id,
		&product.Name,
		&product.DeletedAt,
	)
	if err != nil {
		return entity.Product{}, err
	}
	if product.DeletedAt != nil {
		return entity.Product{}, repoerrs.ErrDataDeleted
	}
	return product, nil
}

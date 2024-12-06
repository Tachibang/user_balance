package repository

import (
	"context"
	"database/sql"
	"log"
	"user_balance/internal/entity"
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
	log.Printf("Создание продукта с именем: %s", name)
	err := r.pg.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при создании продукта: %v", err)
		return 0, err
	}
	log.Printf("Продукт успешно создан, ID: %d", id)
	return id, nil
}

func (r *ProductRepo) GetProduct(ctx context.Context, id int) (entity.Product, error) {
	var product entity.Product
	query := "SELECT * FROM products WHERE id = $1"
	log.Printf("Получение продукта с ID: %d", id)
	err := r.pg.QueryRowContext(ctx, query, id).Scan(
		&product.Id,
		&product.Name,
	)
	if err != nil {
		log.Printf("Ошибка при получении продукта с ID %d: %v", id, err)
		return entity.Product{}, err
	}
	log.Printf("Продукт получен, ID: %d, Имя: %s", product.Id, product.Name)
	return product, nil
}

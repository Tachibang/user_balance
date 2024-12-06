package service

import (
	"context"
	"user_balance/internal/entity"
	"user_balance/internal/repository"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string) (int, error) {
	id, err := s.repo.CreateProduct(ctx, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (entity.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

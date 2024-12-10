package service

import (
	"context"
	"fmt"
	"user_balance/internal/entity"
	"user_balance/internal/repository"

	"github.com/sirupsen/logrus"
)

type ProductService struct {
	repo   repository.Product
	logger *logrus.Logger
}

func NewProductService(repo repository.Product, logger *logrus.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string) (int, error) {
	s.logger.Infof("Создание продукта с именем: %s", name)
	id, err := s.repo.CreateProduct(ctx, name)
	if err != nil {
		err = fmt.Errorf("ошибка при создании продукта с именем %s: %w", name, err)
		s.logger.Error(err)
		return 0, err
	}
	s.logger.Infof("Продукт успешно создан с ID: %d", id)
	return id, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (entity.Product, error) {
	s.logger.Infof("Получение продукта с ID: %d", id)
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		err = fmt.Errorf("ошибка при получении продукта с ID %d: %w", id, err)
		s.logger.Error(err)
		return entity.Product{}, err
	}
	s.logger.Infof("Продукт с ID %d успешно получен: %+v", id, product)
	return product, nil
}

package service

import (
	"context"
	"user_balance/internal/entity"
	"user_balance/internal/repository"

	"github.com/sirupsen/logrus"
)

type ProductService struct {
	repo repository.Product
	log  *logrus.Entry
}

func NewProductService(repo repository.Product, log *logrus.Logger) *ProductService {
	return &ProductService{
		repo: repo,
		log:  log.WithField("service", "product"),
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string) (int, error) {
	id, err := s.repo.CreateProduct(ctx, name)
	s.log.Debugf("CreateProduct, id: %d, name: %s", id, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (entity.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

package service

import (
	"context"
	"user_balance/internal/entity"
	"user_balance/internal/repository"
)

type AccountService struct {
	repo repository.Account
}

func NewAccountService(repo repository.Account) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context) (int, error) {
	id, err := s.repo.CreateAccount(ctx)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id int) (entity.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *AccountService) Deposit(ctx context.Context, id, amount int) (int, int, error) {
	return s.repo.Deposit(ctx, id, amount)
}

func (s *AccountService) Withdraw(ctx context.Context, id, amount int) (int, int, error) {
	return s.repo.Withdraw(ctx, id, amount)
}

func (s *AccountService) Transfer(ctx context.Context, fromID, toID, amount int) (int, int, error) {
	return s.repo.Transfer(ctx, fromID, toID, amount)
}

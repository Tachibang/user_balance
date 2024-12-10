package service

import (
	"context"
	"fmt"
	"user_balance/internal/entity"
	"user_balance/internal/repository"

	"github.com/sirupsen/logrus"
)

type AccountService struct {
	repo   repository.Account
	logger *logrus.Logger
}

func NewAccountService(repo repository.Account, logger *logrus.Logger) *AccountService {
	return &AccountService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context) (int, error) {
	s.logger.Info("Создание нового аккаунта")
	id, err := s.repo.CreateAccount(ctx)
	if err != nil {
		err = fmt.Errorf("ошибка при создании аккаунта: %w", err)
		s.logger.Error(err)
		return 0, err
	}
	s.logger.Infof("Аккаунт успешно создан с ID: %d", id)
	return id, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id int) (entity.Account, error) {
	s.logger.Infof("Получение аккаунта с ID: %d", id)
	account, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		err = fmt.Errorf("ошибка при получении аккаунта с ID %d: %w", id, err)
		s.logger.Error(err)
		return entity.Account{}, err
	}
	s.logger.Infof("Аккаунт с ID %d успешно получен: %+v", id, account)
	return account, nil
}

func (s *AccountService) Deposit(ctx context.Context, id, amount int) (int, int, error) {
	s.logger.Infof("Пополнение аккаунта с ID %d на сумму %d", id, amount)
	balance, totalDeposited, err := s.repo.Deposit(ctx, id, amount)
	if err != nil {
		err = fmt.Errorf("ошибка при пополнении аккаунта с ID %d: %w", id, err)
		s.logger.Error(err)
		return 0, 0, err
	}
	s.logger.Infof("Аккаунт с ID %d успешно пополнен. Баланс: %d, всего пополнений: %d", id, balance, totalDeposited)
	return balance, totalDeposited, nil
}

func (s *AccountService) Withdraw(ctx context.Context, id, amount int) (int, int, error) {
	s.logger.Infof("Снятие с аккаунта с ID %d суммы %d", id, amount)
	balance, totalWithdrawn, err := s.repo.Withdraw(ctx, id, amount)
	if err != nil {
		err = fmt.Errorf("ошибка при снятии с аккаунта с ID %d: %w", id, err)
		s.logger.Error(err)
		return 0, 0, err
	}
	s.logger.Infof("Снятие с аккаунта с ID %d успешно завершено. Баланс: %d, всего снятий: %d", id, balance, totalWithdrawn)
	return balance, totalWithdrawn, nil
}

func (s *AccountService) Transfer(ctx context.Context, fromID, toID, amount int) (int, int, error) {
	s.logger.Infof("Перевод суммы %d с аккаунта %d на аккаунт %d", amount, fromID, toID)
	fromBalance, toBalance, err := s.repo.Transfer(ctx, fromID, toID, amount)
	if err != nil {
		err = fmt.Errorf("ошибка при переводе суммы %d с аккаунта %d на аккаунт %d: %w", amount, fromID, toID, err)
		s.logger.Error(err)
		return 0, 0, err
	}
	s.logger.Infof("Перевод суммы %d с аккаунта %d на аккаунт %d успешно завершен. Баланс отправителя: %d, баланс получателя: %d", amount, fromID, toID, fromBalance, toBalance)
	return fromBalance, toBalance, nil
}

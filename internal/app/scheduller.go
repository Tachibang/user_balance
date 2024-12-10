package app

import (
	"context"
	"time"

	"user_balance/internal/repository"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	cron   *cron.Cron
	logger *logrus.Logger
}

func NewScheduler(logger *logrus.Logger) *Scheduler {
	return &Scheduler{
		cron:   cron.New(),
		logger: logger,
	}
}

func (s *Scheduler) AddJob(spec string, job func()) error {
	_, err := s.cron.AddFunc(spec, job)
	return err
}

func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info("Cron scheduler запущен")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.logger.Info("Cron scheduler остановлен")
}

func (s *Scheduler) GenerateMonthlyReportJob(repo *repository.Repository, producer *Producer, topic string) func() {
	return func() {
		s.logger.Info("Запуск задачи генерации отчета...")
		ctx := context.Background()

		if repo == nil {
			s.logger.Warn("Ошибка: repo равен nil")
			return
		}
		if producer == nil {
			s.logger.Warn("Ошибка: producer равен nil")
			return
		}

		now := time.Now()
		start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		s.logger.Infof("Start Date: %v, End Date: %v", start, end)

		operations, err := repo.GetMonthlyOperations(ctx, start, end)
		if err != nil {
			s.logger.WithError(err).Error("Ошибка получения операций")
			return
		}

		if len(operations) == 0 {
			s.logger.Info("Нет операций для отправки.")
			return
		}

		for _, op := range operations {
			err := producer.SendMessage(topic, op)
			if err != nil {
				s.logger.WithError(err).Error("Ошибка отправки сообщения в Kafka")
			}
		}
	}
}

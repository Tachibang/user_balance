package app

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user_balance/config"
	"user_balance/internal/api"
	"user_balance/internal/api/httpserver"
	"user_balance/internal/repository"
	"user_balance/internal/service"

	_ "github.com/lib/pq"
)

func Run() {
	logger := SetLogrus("info")
	logger.Info("Загрузка конфигурации...")

	cfg, err := config.NewConfig(".env")
	if err != nil {
		logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	logger = SetLogrus(cfg.Level)
	logger.Info("Конфигурация успешно загружена.")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	logger.Info("Подключение к базе данных...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	}()
	logger.Info("Подключение к базе данных успешно установлено.")

	logger.Info("Применение миграций...")
	if err := ApplyMigrations(db, "./migration"); err != nil {
		logger.Fatalf("Ошибка применения миграций: %v", err)
	}
	logger.Info("Миграции успешно применены.")

	logger.Info("Инициализация компонентов приложения...")
	repository := repository.NewRepository(db)
	service := service.NewService(repository, logger)
	router := api.NewRouter(service, logger)
	logger.Info("Компоненты приложения успешно инициализированы.")

	logger.Info("Инициализация Kafka producer...")
	kafkaProducer, err := NewKafkaProducer(cfg.Kafka.Brokers, logger)
	if err != nil {
		logger.Fatalf("Ошибка инициализации Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()
	logger.Info("Kafka producer успешно инициализирован.")

	logger.Info("Инициализация Cron scheduler...")
	scheduler := NewScheduler(logger)
	job := scheduler.GenerateMonthlyReportJob(repository, kafkaProducer, cfg.Kafka.Topic)
	if err := scheduler.AddJob(cfg.Cron.Schedule, job); err != nil {
		logger.Fatalf("Ошибка добавления Cron задачи: %v", err)
	}
	scheduler.Start()
	defer scheduler.Stop()
	logger.Info("Cron scheduler успешно инициализирован.")

	logger.Info("Запуск HTTP сервера...")
	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.Server.Port),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.ShutdownTimeout(15*time.Second),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Warnf("Получен сигнал завершения: %s", s)
	case err = <-httpServer.Notify():
		logger.Errorf("Ошибка HTTP сервера: %v", err)
	}

	logger.Info("Завершение работы HTTP сервера...")
	if err := httpServer.Shutdown(); err != nil {
		logger.Errorf("Ошибка при завершении работы HTTP сервера: %v", err)
	}
	logger.Info("Приложение завершило работу.")
}

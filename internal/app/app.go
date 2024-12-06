package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user_balance/config"
	router "user_balance/internal/api/v1"
	"user_balance/internal/repository"
	"user_balance/internal/service"
	"user_balance/pkg/httpserver"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// TODO: add logger with log levels
func Run() {
	log.Println("Загрузка конфигурации...")
	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	log := configureLogger(cfg.Server.LogLevel)
	log.Info("Конфигурация успешно загружена.")

	// Формируем строку подключения к базе данных
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	log.Info("Подключение к базе данных...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Errorf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	}()
	log.Info("Подключение к базе данных успешно установлено.")

	log.Println("Применение миграций...")
	if err := ApplyMigrations(db, "./migration"); err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
	log.Println("Миграции успешно применены.")

	log.Println("Инициализация компонентов приложения...")
	repository := repository.NewRepository(db)
	service := service.NewService(repository, log)
	router := router.NewV1Router(service)
	log.Println("Компоненты приложения успешно инициализированы.")

	log.Println("Запуск HTTP сервера...")
	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.Server.Port),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.ShutdownTimeout(15*time.Second),
	)

	// Обработка системных сигналов
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Warnf("Получен сигнал завершения: %s", s)
	case err = <-httpServer.Notify():
		log.Errorf("Ошибка HTTP сервера: %v", err)
	}

	log.Println("Завершение работы HTTP сервера...")
	if err := httpServer.Shutdown(); err != nil {
		log.Printf("Ошибка при завершении работы HTTP сервера: %v", err)
	}
	log.Println("Приложение завершило работу.")
}

// zap.Logger, logrus.Logger, std.Logger
func configureLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Fatalf("Ошибка при установке уровня логирования: %v", err)
	}
	logger.SetLevel(lvl)

	return logger
}

/*
DEBUG
INFO
WARN
ERROR
FATAL
*/

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
	"user_balance/internal/api/router"
	"user_balance/internal/repository"
	"user_balance/internal/service"
	"user_balance/pkg/httpserver"

	_ "github.com/lib/pq"
)

func Run() {
	log.Println("Загрузка конфигурации...")
	cfg, err := config.NewConfig(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	log.Println("Конфигурация успешно загружена.")

	// Формируем строку подключения к базе данных
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	log.Println("Подключение к базе данных...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	}()
	log.Println("Подключение к базе данных успешно установлено.")

	log.Println("Применение миграций...")
	if err := ApplyMigrations(db, "./migration"); err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
	log.Println("Миграции успешно применены.")

	log.Println("Инициализация компонентов приложения...")
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	router := router.NewRouter(service)
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
		log.Printf("Получен сигнал завершения: %s", s)
	case err = <-httpServer.Notify():
		log.Printf("Ошибка HTTP сервера: %v", err)
	}

	log.Println("Завершение работы HTTP сервера...")
	if err := httpServer.Shutdown(); err != nil {
		log.Printf("Ошибка при завершении работы HTTP сервера: %v", err)
	}
	log.Println("Приложение завершило работу.")
}

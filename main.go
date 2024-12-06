package main

import (
	"log"
	"user_balance/internal/app"
)

func main() {
	log.Println("Приложение запускается...")

	app.Run()

	log.Println("Приложение успешно завершено.")
}

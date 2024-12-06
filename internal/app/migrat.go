package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ApplyMigrations(db *sql.DB, migrationsPath string) error {
	log.Printf("Чтение директории с миграциями: %s", migrationsPath)
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("Ошибка чтения директории с миграциями: %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(migrationsPath, file.Name())
		log.Printf("Чтение файла миграции: %s", file.Name())
		query, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("Ошибка чтения файла миграции %s: %w", file.Name(), err)
		}

		log.Printf("Применение миграции: %s", file.Name())
		_, err = db.Exec(string(query))
		if err != nil {
			return fmt.Errorf("Ошибка выполнения миграции %s: %w", file.Name(), err)
		}

		log.Printf("Миграция успешно применена: %s", file.Name())
	}
	return nil
}

package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func CheckAndCreateDB() {
	// Получаем значение переменной окружения TODO_DBFILE
	dbFileEnv := os.Getenv("TODO_DBFILE")

	var dbFile string

	if dbFileEnv != "" {
		// Если переменная окружения задана, используем её значение как путь
		dbFile = dbFileEnv
	} else {
		// Если переменная окружения не задана, формируем путь по умолчанию
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		// Формируем путь к файлу БД
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	// Проверяем существование файла
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		// Файл не существует — создаем новый
		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE scheduler (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		date DATE NOT NULL,
    		title VARCHAR(255) NOT NULL,
    		comment TEXT,
    		repeat VARCHAR(128),
    		CONSTRAINT unique_date_title UNIQUE (date, title)
		)`)
		if err != nil {
			log.Fatalf("Ошибка при создании таблицы: %v", err)
		}
		log.Println("База данных создана успешно.")
	} else {
		log.Println("Файл базы данных уже существует.")
	}
}

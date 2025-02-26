package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func CheckAndCreateDB() *sql.DB {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении текущей директории: %v", err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "go_final_project", "scheduler.db")
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = envFile
	}
	_, err = os.Stat(filepath.Join(dbFile))

	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			log.Fatalf("Ошибка при открытии базы данных: %v", err)
		}

		log.Println("Создание новой базы данных и таблицы...")
		createTableSql := `CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT "",
            title VARCHAR(128) NOT NULL DEFAULT "",
            comment TEXT NOT NULL DEFAULT "",
            repeat VARCHAR(128) NOT NULL DEFAULT ""
            );
            CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`

		_, err = db.Exec(createTableSql)
		if err != nil {
			log.Fatalf("Ошибка при создании таблицы: %v", err)
		}
		log.Println("База данных и таблица успешно созданы.")
		return db
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatalf("Ошибка при открытии базы данных: %v", err)
	}
	log.Println("База данных успешно открыта.")
	return db
}

package pkg

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func CheckAndCreateDB() {
	// Получаем значение переменной окружения TODO_DBFILE
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile != "" {
		log.Printf("Переменная окружения TODO_DBFILE задана: %s", dbFile)
	} else {
		log.Println("Переменная окружения TODO_DBFILE не задана, используем путь по умолчанию")
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	log.Printf("Путь к файлу базы данных: %s", dbFile)

	// Проверяем, существует ли файл БД
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		// Файл не существует — создаём базу данных
		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Создаём таблицу scheduler с нужными полями
		createTableQuery := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(length(repeat) <= 128),
			CONSTRAINT unique_date_title UNIQUE (date, title)
		);
		`
		if _, err = db.Exec(createTableQuery); err != nil {
			log.Fatalf("Ошибка при создании таблицы: %v", err)
		}

		// Создаём индекс по полю date для сортировки задач
		createIndexQuery := `CREATE INDEX idx_scheduler_date ON scheduler(date);`
		if _, err = db.Exec(createIndexQuery); err != nil {
			log.Fatalf("Ошибка при создании индекса: %v", err)
		}

		log.Println("База данных создана успешно.")
	} else {
		log.Println("Файл базы данных уже существует.")
	}
}

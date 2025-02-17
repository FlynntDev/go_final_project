package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/flynntdev/go_final_project/db"
)

// webDir -путь к директории с файлами, которые сервер будет раздавать
const webDir = "./web"

func main() {

	db.CheckAndCreateDB()
	// Получаем порт из переменной окружения TODO_PORT, если не задано, используем 7540
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	fmt.Printf("Запуск сервера на порту %s\n", port)

	// Создаём обработчик, который раздаёт файлы из указанной директории
	fileServer := http.FileServer(http.Dir(webDir))

	// Используем http.StripPrefix, чтобы корректно обрабатывать запросы к файлам
	http.Handle("/", http.StripPrefix("/", fileServer))

	// Запускаем сервер на указанном порту
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err) // В случае ошибки выводим её и завершаем программу
	}

	fmt.Println("Завершаем работу")
}

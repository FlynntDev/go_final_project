package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/flynntdev/go_final_project/api"
)

// webDir -путь к директории с файлами, которые сервер будет раздавать
const webDir = "./web"

func ServerStart() {
	// Получаем порт
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fmt.Printf("Запуск сервера на порту %s\n", port)

	// Регистрируем обработчик API
	http.HandleFunc("/api/nextdate", NextDateHandler)

	// Регистрируем обработчик для раздачи статических файлов
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", http.StripPrefix("/", fileServer))

	// Запускаем сервер
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса через FormValue
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверяем, что параметры присутствуют
	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "отсутствуют необходимые параметры", http.StatusBadRequest)
		return
	}

	// Парсим дату "now"
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "некорректная дата now", http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату с использованием переданных параметров
	nextDate, err := api.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем просто дату в ответ
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

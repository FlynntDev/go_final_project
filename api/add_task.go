package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HandleAddTask(w http.ResponseWriter, r *http.Request) {
	var task struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Получение текущей даты
	now := time.Now()
	currentDate := now.Format("20060102")

	// Проверка и обработка даты
	if task.Date == "" {
		task.Date = currentDate
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		if parsedDate.Before(now) {
			if task.Repeat == "" {
				task.Date = currentDate
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				task.Date = nextDate
			}
		}
	}

	// Печать полученных данных в консоль
	fmt.Printf("Получена задача: %+v\n", task)

	// Ответ для клиента
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("TESTTTTT"))
}

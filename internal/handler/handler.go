package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/flynntdev/go_final_project/api"
	"github.com/flynntdev/go_final_project/config"
	"github.com/flynntdev/go_final_project/internal/storage"
	"github.com/flynntdev/go_final_project/internal/task"
)

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func HandlerPostGetPutTask(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t task.Task
		switch {
		case r.Method == http.MethodPost:
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				log.Printf("Ошибка декодирования JSON: %v", err)
				http.Error(w, `{"error":"Ошибка преобразования JSON"}`, http.StatusBadRequest)
				return
			}
			id, err := store.PostTask(t)
			if err != nil {
				log.Printf("Ошибка создания задачи: %v", err)
				http.Error(w, `{"error":"Ошибка создания задачи"}`, http.StatusBadRequest)
				return
			}

			resp := Response{ID: id}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Ошибка кодирования ответа в формате JSON: %v", err)
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case r.Method == http.MethodGet:
			id := r.URL.Query().Get("id")
			task, err := store.GetTask(id)
			if err != nil {
				log.Printf("Ошибка получения задачи: %v", err)
				http.Error(w, `{"error":"Ошибка получения задачи"}`, http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(task); err != nil {
				log.Printf("Ошибка кодирования ответа в формате JSON: %v", err)
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case r.Method == http.MethodPut:
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				log.Printf("Ошибка декодирования JSON: %v", err)
				http.Error(w, `{"error":"Ошибка преобразования JSON"}`, http.StatusBadRequest)
				return
			}
			err = store.PutTask(t)
			if err != nil {
				log.Printf("Ошибка обновления задачи: %v", err)
				http.Error(w, `{"error":"Ошибка обновления задачи"}`, http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
				log.Printf("Ошибка кодирования ответа в формате JSON: %v", err)
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}

		case r.Method == http.MethodDelete:
			id := r.URL.Query().Get("id")
			err := store.DeleteTask(id)
			if err != nil {
				log.Printf("Ошибка удаления задачи: %v", err)
				http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
				log.Printf("Ошибка кодирования ответа в формате JSON: %v", err)
				http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
				return
			}
		}
	}
}

func HandlerGetTasks(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		tasks, err := store.SearchTask(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := map[string][]task.Task{
			"tasks": tasks,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

func HandlerNextDate(w http.ResponseWriter, r *http.Request) {

	strnow := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	strRepeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(config.Layout, strnow)
	if err != nil {
		log.Println("Ошибка парсинга времени:", err)
		http.Error(w, `{"error":"Неверный формат времени"}`, http.StatusBadRequest)
		return
	}

	nextdate, err := api.NextDate(now, date, strRepeat)
	if err != nil {
		log.Println("Ошибка получения следующей даты:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(nextdate))
	if err != nil {
		log.Println("Ошибка записи ответа:", err)
		http.Error(w, `{"error":"Ошибка записи ответа"}`, http.StatusInternalServerError)
		return
	}
}

func HandlerDone(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		err := store.DoneTask(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
			http.Error(w, `{"error":"Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}
	}
}

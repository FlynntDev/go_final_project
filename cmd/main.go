package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/flynntdev/go_final_project/api"
	"github.com/flynntdev/go_final_project/config"
	"github.com/flynntdev/go_final_project/db"
	"github.com/flynntdev/go_final_project/internal/hendler"
	"github.com/flynntdev/go_final_project/internal/storage"

	"github.com/go-chi/chi"
)

var Password string

func main() {
	env := config.GetEnv()
	Password = os.Getenv("TODO_PASSWORD")
	fmt.Println("Приложение запущено на порту", env.Port)

	dataBase := db.CheckAndCreateDB()
	defer dataBase.Close()
	store := storage.NewStore(dataBase)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", hendler.HandlerNextDate)
	r.Post("/api/task", api.Authorization(hendler.HandlerPostGetPutTask(store)))
	r.Get("/api/tasks", api.Authorization(hendler.HandlerGetTasks(store)))
	r.Get("/api/task", hendler.HandlerPostGetPutTask(store))
	r.Put("/api/task", hendler.HandlerPostGetPutTask(store))
	r.Post("/api/task/done", api.Authorization(hendler.HandlerDone(store)))
	r.Delete("/api/task", hendler.HandlerPostGetPutTask(store))
	r.Post("/api/signin", api.SigninHandler)

	err := http.ListenAndServe(":"+env.Port, r)
	if err != nil {
		fmt.Println("ошибка запуска сервера:", err)
	}
}

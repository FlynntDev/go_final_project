package main

import (
	"github.com/flynntdev/go_final_project/internal/pkg/db"
	"github.com/flynntdev/go_final_project/server"
)

func main() {
	// Проверка и создание БД
	db.CheckAndCreateDB()
	// Старт сервера
	server.ServerStart()

}

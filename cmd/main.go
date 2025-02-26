package main

import (
	"github.com/flynntdev/go_final_project/internal/pkg"
	"github.com/flynntdev/go_final_project/server"
)

func main() {
	// Проверка и создание БД
	pkg.CheckAndCreateDB()
	// Старт сервера
	server.ServerStart()

}

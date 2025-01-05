package main

import (
	"forum/backend/database"
	"forum/backend/handlers"
	"forum/backend/server"
)

func main() {
	database.CreateDatabaseIfNotExists()

	handlers.ImportHandlers()

	server.StartServer()
}

package server

import (
	"fmt"
	"net/http"

	"forum/backend/controllers/login"
)

const port = ":8080"

func StartServer() {
	login.LoadEnv()

	fmt.Printf("Server is running on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}

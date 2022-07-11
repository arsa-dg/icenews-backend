package main

import (
	"context"
	"icenews/backend/config"
	"icenews/backend/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	DB := config.ConnectDB()

	defer DB.Close(context.Background())

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("IceNews Backend"))
	})

	router.Mount("/auth", routes.AuthRoute(DB))

	http.ListenAndServe(":8080", router)
}

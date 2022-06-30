package main

import (
	"context"
	"fmt"
	"icenews/backend/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func main() {
	// hardcoded dburl sementara
	DBUrl := "postgres://postgres:password@localhost:5432/icenews"

	DB, err := pgx.Connect(context.Background(), DBUrl)
	if err != nil {
		fmt.Println("Error while connecting to database!")
	}
	defer DB.Close(context.Background())

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("IceNews Backend"))
	})

	router.Mount("/auth", routes.AuthRoute(DB))

	http.ListenAndServe(":8080", router)
}

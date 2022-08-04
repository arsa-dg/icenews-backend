package main

import (
	"context"
	"icenews/backend/config"
	"icenews/backend/repository"
	"icenews/backend/routes"
	"icenews/backend/service"
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

	userRepository := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepository)

	newsRepository := repository.NewNewsRepository(DB)
	newsService := service.NewNewsService(newsRepository)

	router.Mount("/auth", routes.AuthRoute(userService))
	router.Mount("/me", routes.MeRoute(userService))
	router.Mount("/news", routes.NewsRoute(newsService))

	http.ListenAndServe(":8080", router)
}

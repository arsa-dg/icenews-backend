package main

import (
	"icenews/backend/config"
	"icenews/backend/middleware"
	"icenews/backend/repository"
	"icenews/backend/routes"
	"icenews/backend/service"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load(".env")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Server is listening on port 8080")

	DB := config.ConnectDB()

	defer DB.Close()

	router := chi.NewRouter()
	router.Use(middleware.RequestLogger)

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

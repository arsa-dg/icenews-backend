package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"
	"icenews/backend/repository"
	"icenews/backend/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func AuthRoute(DB *pgx.Conn) chi.Router {
	userRepository := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(userService)

	r := chi.NewRouter()

	r.Post("/login", authHandler.Login)
	r.Post("/register", authHandler.Register)
	r.With(middleware.MiddlewareAuth).Get("/token", authHandler.Token)

	return r
}

package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"
	"icenews/backend/repository"
	"icenews/backend/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func MeRoute(DB *pgx.Conn) chi.Router {
	userRepository := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepository)
	meHandler := handler.NewMeHandler(userService)

	r := chi.NewRouter()

	r.With(middleware.MiddlewareAuth).Get("/profile", meHandler.Profile)

	return r
}

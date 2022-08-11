package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"
	"icenews/backend/service"

	"github.com/go-chi/chi/v5"
)

func AuthRoute(s service.UserServiceInterface) chi.Router {
	authHandler := handler.NewAuthHandler(s)

	r := chi.NewRouter()

	r.With(middleware.TypeJSON).Post("/login", authHandler.Login)
	r.With(middleware.TypeJSON).Post("/register", authHandler.Register)
	r.With(middleware.JWT).Get("/token", authHandler.Token)

	return r
}

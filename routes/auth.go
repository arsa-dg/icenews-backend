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

	r.Post("/login", authHandler.Login)
	r.Post("/register", authHandler.Register)
	r.With(middleware.JWT).Get("/token", authHandler.Token)

	return r
}

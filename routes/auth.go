package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func AuthRoute(DB *pgx.Conn) chi.Router {
	authHandler := handler.NewAuthHandler(DB)

	r := chi.NewRouter()

	r.Post("/login", authHandler.Login)
	r.Post("/register", authHandler.Register)
	r.With(middleware.MiddlewareAuth).Get("/token", authHandler.Token)

	return r
}

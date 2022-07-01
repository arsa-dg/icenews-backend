package routes

import (
	"icenews/backend/handler/auth"
	"icenews/backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func AuthRoute(DB *pgx.Conn) chi.Router {
	authHandler := auth.NewAuthHandler(DB)

	r := chi.NewRouter()

	r.Post("/login", authHandler.Login)
	r.With(middleware.MiddlewareAuth).Get("/token", authHandler.Token)

	// r.Group(func(r chi.Router) {
	// 	r.Use(middleware.MiddlewareAuth)
	// 	r.Get("/token", authHandler.Token)
	// })

	return r
}

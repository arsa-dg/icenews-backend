package routes

import (
	"icenews/backend/handler/auth"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func AuthRoute(DB *pgx.Conn) chi.Router {
	authHandler := auth.NewAuthHandler(DB)

	r := chi.NewRouter()

	r.Post("/login", authHandler.Login)

	return r
}

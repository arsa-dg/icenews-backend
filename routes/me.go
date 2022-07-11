package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func MeRoute(DB *pgx.Conn) chi.Router {
	meHandler := handler.NewMeHandler(DB)

	r := chi.NewRouter()

	r.With(middleware.MiddlewareAuth).Get("/profile", meHandler.Profile)

	return r
}

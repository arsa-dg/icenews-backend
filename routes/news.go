package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func NewsRoute(DB *pgx.Conn) chi.Router {
	authHandler := handler.NewNewsHandler(DB)

	r := chi.NewRouter()

	r.With(middleware.MiddlewareAuth).Get("/", authHandler.GetAll)

	return r
}

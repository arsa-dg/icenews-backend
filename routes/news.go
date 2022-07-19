package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func NewsRoute(DB *pgx.Conn) chi.Router {
	newsHandler := handler.NewNewsHandler(DB)

	r := chi.NewRouter()

	r.With(middleware.MiddlewareAuth).Get("/", newsHandler.GetAll)
	r.With(middleware.MiddlewareAuth).Get("/{id:^[0-9]+}", newsHandler.GetDetail)

	return r
}

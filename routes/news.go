package routes

import (
	"icenews/backend/handler"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func NewsRoute(DB *pgx.Conn) chi.Router {
	authHandler := handler.NewNewsHandler(DB)

	r := chi.NewRouter()

	r.Get("/", authHandler.GetAll)

	return r
}

package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"
	"icenews/backend/repository"
	"icenews/backend/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

func NewsRoute(DB *pgx.Conn) chi.Router {
	newsRepository := repository.NewNewsRepository(DB)
	newsService := service.NewNewsService(newsRepository)
	newsHandler := handler.NewNewsHandler(newsService)

	r := chi.NewRouter()

	r.Use(middleware.MiddlewareAuth)

	r.Get("/", newsHandler.GetAll)
	r.Get("/{id:^[0-9]+}", newsHandler.GetDetail)
	r.Get("/category", newsHandler.NewsCategory)
	r.Post("/{id:^[0-9]+}/comment", newsHandler.AddComment)
	r.Get("/{id:^[0-9]+}/comment", newsHandler.CommentList)

	return r
}

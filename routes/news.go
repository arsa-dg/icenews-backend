package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"
	"icenews/backend/service"

	"github.com/go-chi/chi/v5"
)

func NewsRoute(s service.NewsServiceInterface) chi.Router {
	newsHandler := handler.NewNewsHandler(s)

	r := chi.NewRouter()

	r.Use(middleware.JWT)

	r.Get("/", newsHandler.GetAll)
	r.Get("/{id:^[0-9]+}", newsHandler.GetDetail)
	r.Get("/category", newsHandler.NewsCategory)
	r.With(middleware.TypeJSON).Post("/{id:^[0-9]+}/comment", newsHandler.AddComment)
	r.Get("/{id:^[0-9]+}/comment", newsHandler.CommentList)

	return r
}

package routes

import (
	"icenews/backend/handler"
	"icenews/backend/middleware"
	"icenews/backend/service"

	"github.com/go-chi/chi/v5"
)

func MeRoute(s service.UserServiceInterface) chi.Router {
	meHandler := handler.NewMeHandler(s)

	r := chi.NewRouter()

	r.With(middleware.JWT).Get("/profile", meHandler.Profile)

	return r
}

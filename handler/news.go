package handler

import (
	"icenews/backend/helper"
	"icenews/backend/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
)

type NewsHandler struct {
	NewsService service.NewsService
}

func NewNewsHandler(DB *pgx.Conn) NewsHandler {
	return NewsHandler{service.NewNewsService(DB)}
}

func (h NewsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	response, statusCode := h.NewsService.GetAllLogic(r.URL.Query())

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)

		return
	}

	helper.ResponseOK(w, response)
}

func (h NewsHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	newsId := chi.URLParam(r, "id")

	response, statusCode := h.NewsService.GetDetailLogic(newsId)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)

		return
	}

	helper.ResponseOK(w, response)
}

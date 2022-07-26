package handler

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NewsHandler struct {
	NewsService service.NewsService
}

func NewNewsHandler(s service.NewsService) NewsHandler {
	return NewsHandler{s}
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

func (h NewsHandler) NewsCategory(w http.ResponseWriter, r *http.Request) {
	response, statusCode := h.NewsService.NewsCategoryLogic()

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)

		return
	}

	helper.ResponseOK(w, response)
}

func (h NewsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	newsId := chi.URLParam(r, "id")
	var field interfaces.CommentRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Wrong Request Format",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)

		return
	}

	userId := r.Context().Value("user_id").(string)
	userIdUUID, err := uuid.Parse(userId)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusInternalServerError, res)

		return
	}

	response, statusCode := h.NewsService.AddCommentLogic(field, newsId, userIdUUID)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)

		return
	}

	helper.ResponseOK(w, response)
}

func (h NewsHandler) CommentList(w http.ResponseWriter, r *http.Request) {
	newsId := chi.URLParam(r, "id")

	response, statusCode := h.NewsService.CommentListLogic(newsId)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)

		return
	}

	helper.ResponseOK(w, response)
}

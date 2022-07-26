package handler

import (
	"icenews/backend/helper"
	"icenews/backend/model"
	"icenews/backend/service"
	"net/http"

	"github.com/google/uuid"
)

type MeHandlerInterface interface {
	Profile(w http.ResponseWriter, r *http.Request)
}

type MeHandler struct {
	UserService service.UserServiceInterface
}

func NewMeHandler(s service.UserServiceInterface) MeHandler {
	return MeHandler{s}
}

func (h MeHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	userIdUUID, err := uuid.Parse(userId)

	if err != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusInternalServerError, res)

		return
	}

	res, statusCode := h.UserService.ProfileLogic(userIdUUID)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, res)

		return
	}

	helper.ResponseOK(w, res)
}

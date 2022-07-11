package handler

import (
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type MeHandler struct {
	UserService service.UserService
}

func NewMeHandler(DB *pgx.Conn) MeHandler {
	return MeHandler{service.NewUserService(DB)}
}

func (h MeHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	userIdUUID, err := uuid.Parse(userId)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusInternalServerError, res)
	}

	res, statusCode := h.UserService.ProfileLogic(userIdUUID)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, res)
	}

	helper.ResponseOK(w, res)
}

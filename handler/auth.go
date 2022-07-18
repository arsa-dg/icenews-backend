package handler

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/service"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type AuthHandler struct {
	UserService service.UserService
}

func NewAuthHandler(DB *pgx.Conn) AuthHandler {
	return AuthHandler{service.NewUserService(DB)}
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var field interfaces.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Wrong Request Format",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)

		return
	}

	response, statusCode := h.UserService.LoginLogic(field)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)
		return
	}

	helper.ResponseOK(w, response)
}

func (h AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	token, expiresAt, err := helper.CreateJWT(userId)

	// bad request (400)
	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusInternalServerError, res)

		return
	}

	res := interfaces.AuthLoginResponse{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	helper.ResponseOK(w, res)
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var field interfaces.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Wrong Request Format",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)

		return
	}

	response, statusCode := h.UserService.RegisterLogic(field)

	if statusCode != http.StatusOK {
		helper.ResponseError(w, statusCode, response)
		return
	}

	helper.ResponseOK(w, response)
}

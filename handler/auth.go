package handler

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/model"
	"icenews/backend/service"
	"net/http"
)

type AuthHandlerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
	Token(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	UserService service.UserServiceInterface
}

func NewAuthHandler(s service.UserServiceInterface) AuthHandler {
	return AuthHandler{s}
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var field model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := model.ResponseBadRequest{
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
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusInternalServerError, res)

		return
	}

	res := model.AuthLoginResponse{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	helper.ResponseOK(w, res)
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var field model.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := model.ResponseBadRequest{
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

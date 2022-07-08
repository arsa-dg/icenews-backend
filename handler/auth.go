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
	DB *pgx.Conn
}

func NewAuthHandler(DB *pgx.Conn) AuthHandler {
	return AuthHandler{DB}
}

func (Handler AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var field interfaces.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Wrong Request Format",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)

		return
	}

	userService := service.NewUserService(Handler.DB)

	response, statusCode := userService.LoginLogic(field)

	if statusCode == http.StatusOK {
		helper.ResponseOK(w, response)
	} else {
		helper.ResponseError(w, statusCode, response)
	}
}

func (AH AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	token, expiresAt, err := helper.CreateJWT(userId)

	// bad request (400)
	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)
	}

	res := interfaces.AuthResponseOK{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	helper.ResponseOK(w, res)
}

func (Handler AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var field interfaces.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&field)

	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Wrong Request Format",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)

		return
	}

	userService := service.NewUserService(Handler.DB)

	response, statusCode := userService.RegisterLogic(field)

	if statusCode == http.StatusOK {
		helper.ResponseOK(w, response)
	} else {
		helper.ResponseError(w, statusCode, response)
	}
}

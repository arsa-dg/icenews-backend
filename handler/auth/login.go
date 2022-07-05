package auth

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/service"
	"net/http"
)

func (Handler AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// to do
	// validate hash salt?
	// validate username password (min 8 char, etc)?
	// response 401?

	var field interfaces.LoginRequest
	json.NewDecoder(r.Body).Decode(&field)

	userService := service.NewUserService(Handler.DB)

	response, statusCode := userService.LoginLogic(field)

	if statusCode == http.StatusOK {
		helper.ResponseOK(w, response)
	} else {
		helper.ResponseError(w, statusCode, response)
	}
}

package auth

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/service"
	"net/http"
)

func (Handler AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

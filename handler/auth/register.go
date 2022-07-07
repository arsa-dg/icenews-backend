package auth

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/service"
	"net/http"
)

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

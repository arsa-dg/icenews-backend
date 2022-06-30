package auth

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/repository"
	"net/http"
)

type LoginField struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseOK struct {
	Token      string `json:"token"`
	Scheme     string `json:"scheme"`
	Expires_at string `json:"expires_at"`
}

func (AH AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// to do
	// validate hash salt?
	// validate username/password empty?
	// jwt
	// response 400, 401, 422

	var field LoginField
	json.NewDecoder(r.Body).Decode(&field)

	userRepository := repository.NewUserRepository(AH.DB)

	user := userRepository.SelectByUsername(field.Username)

	if user.Password == field.Password {
		res := responseOK{
			Token:      "tes", // to be implemented
			Scheme:     "Bearer",
			Expires_at: "tes", // to be implemented
		}

		helper.ResponseOK(w, res)
	}
}

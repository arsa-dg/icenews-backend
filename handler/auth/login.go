package auth

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponseOK struct {
	Token      string `json:"token"`
	Scheme     string `json:"scheme"`
	Expires_at string `json:"expires_at"`
}

func (AH AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// to do
	// validate hash salt?
	// validate username password (min 8 char, etc)?
	// response 401?

	var field LoginRequest
	json.NewDecoder(r.Body).Decode(&field)

	// field empty (validation error)
	if field.Username == "" || field.Password == "" {
		res := interfaces.ResponseValidationFailed{
			Message: "Field(s) is(are) missing",
		}

		var emptyFields []interfaces.FieldError

		if field.Username == "" {
			toAdd := interfaces.FieldError{
				Name:  "username",
				Error: "username is missing",
			}

			emptyFields = append(emptyFields, toAdd)
		}

		if field.Password == "" {
			toAdd := interfaces.FieldError{
				Name:  "password",
				Error: "password is missing",
			}

			emptyFields = append(emptyFields, toAdd)
		}

		res.Field = emptyFields

		helper.ResponseError(w, http.StatusUnprocessableEntity, res)

	} else {
		userRepository := repository.NewUserRepository(AH.DB)

		user := userRepository.SelectByUsername(field.Username)

		// ok
		if user.Password == field.Password {
			token, expiresAt := helper.CreateJWT(user.Id)

			res := AuthResponseOK{
				Token:      token,
				Scheme:     "Bearer",
				Expires_at: expiresAt,
			}

			helper.ResponseOK(w, res)

			// wrong password (bad request)
		} else {
			res := interfaces.ResponseBadRequest{
				Message: "Wrong Password",
			}

			helper.ResponseError(w, http.StatusBadRequest, res)
		}
	}
}

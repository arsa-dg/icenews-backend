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

type responseBadRequest struct {
	Message string `json:"message"`
}

type fieldError struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

type responseValidationFailed struct {
	Message string       `json:"message"`
	Field   []fieldError `json:"field"`
}

func (AH AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// to do
	// validate hash salt?
	// validate username password (min 8 char, etc)?
	// response 401?

	var field LoginField
	json.NewDecoder(r.Body).Decode(&field)

	// field empty (validation error)
	if field.Username == "" || field.Password == "" {
		res := responseValidationFailed{
			Message: "Field(s) is(are) missing",
		}

		var emptyFields []fieldError

		if field.Username == "" {
			toAdd := fieldError{
				Name:  "username",
				Error: "username is missing",
			}

			emptyFields = append(emptyFields, toAdd)
		}

		if field.Password == "" {
			toAdd := fieldError{
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

			res := responseOK{
				Token:      token,
				Scheme:     "Bearer",
				Expires_at: expiresAt,
			}

			helper.ResponseOK(w, res)

			// wrong password (bad request)
		} else {
			res := responseBadRequest{
				Message: "Wrong Password",
			}

			helper.ResponseError(w, http.StatusBadRequest, res)
		}
	}
}

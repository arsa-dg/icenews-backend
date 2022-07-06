package service

import (
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *pgx.Conn
}

func NewUserService(DB *pgx.Conn) UserService {
	return UserService{DB}
}

func (Service UserService) LoginLogic(request interfaces.LoginRequest) (interface{}, int) {
	// field empty (validation error 422)
	if helper.IsEmptyStrings(request.Username, request.Password) {
		res := interfaces.ResponseValidationFailed{
			Message: "Field(s) is(are) missing",
		}

		var emptyFields []interfaces.FieldError

		if helper.IsEmptyStrings(request.Username) {
			toAdd := interfaces.FieldError{
				Name:  "username",
				Error: "username is missing",
			}

			emptyFields = append(emptyFields, toAdd)
		}

		if helper.IsEmptyStrings(request.Password) {
			toAdd := interfaces.FieldError{
				Name:  "password",
				Error: "password is missing",
			}

			emptyFields = append(emptyFields, toAdd)
		}

		res.Field = emptyFields

		return res, http.StatusUnprocessableEntity
	} else {
		userRepository := repository.NewUserRepository(Service.DB)
		user := userRepository.SelectByUsername(request.Username)

		isPassErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

		// wrong password (bad request 400)
		if isPassErr == bcrypt.ErrMismatchedHashAndPassword {
			res := interfaces.ResponseBadRequest{
				Message: "Wrong Password",
			}

			return res, http.StatusBadRequest
		} else { // ok
			token, expiresAt := helper.CreateJWT(user.Id)

			res := interfaces.AuthResponseOK{
				Token:      token,
				Scheme:     "Bearer",
				Expires_at: expiresAt,
			}

			return res, http.StatusOK
		}
	}
}

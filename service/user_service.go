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
	}

	userRepository := repository.NewUserRepository(Service.DB)
	user, err := userRepository.SelectByUsername(request.Username)

	// user not found (invalid credentials 401)
	if err == pgx.ErrNoRows {
		res := interfaces.ResponseUnauthorized{
			Message: "User Not Found",
		}

		return res, http.StatusUnauthorized
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	// wrong password (invalid credentials 401)
	if errPass == bcrypt.ErrMismatchedHashAndPassword {
		res := interfaces.ResponseUnauthorized{
			Message: "Wrong Password",
		}

		return res, http.StatusUnauthorized
	}

	token, expiresAt := helper.CreateJWT(user.Id.String())

	res := interfaces.AuthResponseOK{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	return res, http.StatusOK
}

func (Service UserService) RegisterLogic(request interfaces.RegisterRequest) (interface{}, int) {
	if helper.IsEmptyStrings(request.Username, request.Password, request.Name, request.Bio, request.Web, request.Picture) {
		res := interfaces.ResponseValidationFailed{
			Message: "Field(s) is(are) missing",
		}

		return res, http.StatusUnprocessableEntity
	}

	userRepository := repository.NewUserRepository(Service.DB)
	_, err := userRepository.SelectByUsername(request.Username)

	if err != pgx.ErrNoRows {
		res := interfaces.ResponseBadRequest{
			Message: "Username Is Not Available",
		}

		return res, http.StatusBadRequest
	}

	hashPass, errGenerate := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if errGenerate != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Something Is Wrong",
		}

		return res, http.StatusBadRequest
	}

	errInsert := userRepository.Insert(request.Username, string(hashPass), request.Name, request.Bio, request.Web, request.Picture)

	if errInsert != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Something Is Wrong",
		}

		return res, http.StatusBadRequest
	}

	res := interfaces.ResponseOK{
		Message: "Register Success",
	}

	return res, http.StatusOK
}

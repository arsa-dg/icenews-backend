package service

import (
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB        *pgx.Conn
	Validator *validator.Validate
}

func NewUserService(DB *pgx.Conn) UserService {
	return UserService{DB, validator.New()}
}

func (Service UserService) LoginLogic(request interfaces.LoginRequest) (interface{}, int) {
	// field empty (validation error 422)
	errValidateRes, errValidateStatus := helper.RequestValidation(Service.Validator, request)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	userRepository := repository.NewUserRepository(Service.DB)
	user, errSelect := userRepository.SelectByUsername(request.Username)

	// user not found (invalid credentials 401)
	if errSelect == pgx.ErrNoRows {
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

	token, expiresAt, errGenerate := helper.CreateJWT(user.Id.String())

	// bad request (400)
	if errGenerate != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Something Is Wrong",
		}

		return res, http.StatusBadRequest
	}

	// OK (200)
	res := interfaces.AuthResponseOK{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	return res, http.StatusOK
}

func (Service UserService) RegisterLogic(request interfaces.RegisterRequest) (interface{}, int) {
	errValidateRes, errValidateStatus := helper.RequestValidation(Service.Validator, request)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
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

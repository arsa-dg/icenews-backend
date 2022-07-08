package service

import (
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (s UserService) LoginLogic(request interfaces.LoginRequest) (interface{}, int) {
	// field empty (validation error 422)
	errValidateRes, errValidateStatus := helper.RequestValidation(s.Validator, request)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	userRepository := repository.NewUserRepository(s.DB)
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
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	// OK (200)
	res := interfaces.AuthLoginResponse{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	return res, http.StatusOK
}

func (s UserService) RegisterLogic(request interfaces.RegisterRequest) (interface{}, int) {
	errValidateRes, errValidateStatus := helper.RequestValidation(s.Validator, request)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	userRepository := repository.NewUserRepository(s.DB)
	_, err := userRepository.SelectByUsername(request.Username)

	if err != pgx.ErrNoRows {
		res := interfaces.ResponseBadRequest{
			Message: "Username Is Not Available",
		}

		return res, http.StatusBadRequest
	}

	hashPass, errGenerate := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if errGenerate != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	id := uuid.New()
	newUser := interfaces.User{
		Id:       id,
		Username: request.Username,
		Password: string(hashPass),
		Name:     request.Name,
		Bio:      request.Bio,
		Web:      request.Web,
		Picture:  request.Picture,
	}

	errInsert := userRepository.Insert(newUser)

	if errInsert != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := interfaces.ResponseOK{
		Message: "Register Success",
	}

	return res, http.StatusOK
}

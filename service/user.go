package service

import (
	"icenews/backend/helper"
	"icenews/backend/model"
	"icenews/backend/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	LoginLogic(request model.LoginRequest) (interface{}, int)
	TokenLogic(id string) (interface{}, int)
	RegisterLogic(request model.RegisterRequest) (interface{}, int)
	ProfileLogic(id uuid.UUID) (interface{}, int)
}

type UserService struct {
	Validator      *validator.Validate
	UserRepository repository.UserRepositoryInterface
}

func NewUserService(r repository.UserRepositoryInterface) UserService {
	return UserService{validator.New(), r}
}

func (s UserService) LoginLogic(request model.LoginRequest) (interface{}, int) {
	// field empty (validation error 422)
	errValidateRes, errValidateStatus := helper.RequestValidation(s.Validator, request)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	user, errSelect := s.UserRepository.SelectByUsername(request.Username)

	// user not found (invalid credentials 401)
	if errSelect != nil || user.Username == "" {
		res := model.ResponseUnauthorized{
			Message: "User Not Found",
		}

		return res, http.StatusUnauthorized
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	// wrong password (invalid credentials 401)
	if errPass == bcrypt.ErrMismatchedHashAndPassword {
		res := model.ResponseUnauthorized{
			Message: "Wrong Password",
		}

		return res, http.StatusUnauthorized
	}

	token, expiresAt, errGenerate := helper.CreateJWT(user.Id.String())

	// internal server error (500)
	if errGenerate != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	// OK (200)
	res := model.AuthLoginResponse{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	return res, http.StatusOK
}

func (s UserService) TokenLogic(id string) (interface{}, int) {
	userIdUUID, errParse := uuid.Parse(id)

	if errParse != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	user, errSelect := s.UserRepository.SelectById(userIdUUID)

	if errSelect != nil || user.Username == "" {
		res := model.ResponseNotFound{
			Message: "User Not Found",
		}

		return res, http.StatusNotFound
	}

	token, expiresAt, errGenerate := helper.CreateJWT(id)

	// internal server error (500)
	if errGenerate != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	// OK (200)
	res := model.AuthLoginResponse{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	return res, http.StatusOK
}

func (s UserService) RegisterLogic(request model.RegisterRequest) (interface{}, int) {
	errValidateRes, errValidateStatus := helper.RequestValidation(s.Validator, request)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	user, errSelect := s.UserRepository.SelectByUsername(request.Username)

	if errSelect == nil || user.Username == request.Username {
		res := model.ResponseBadRequest{
			Message: "Username Is Not Available",
		}

		return res, http.StatusBadRequest
	}

	hashPass, errGenerate := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if errGenerate != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	id := uuid.New()
	newUser := model.User{
		Id:       id,
		Username: request.Username,
		Password: string(hashPass),
		Name:     request.Name,
		Bio:      request.Bio,
		Web:      request.Web,
		Picture:  request.Picture,
	}

	errInsert := s.UserRepository.Insert(newUser)

	if errInsert != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := model.ResponseOK{
		Message: "Register Success",
	}

	return res, http.StatusOK
}

func (s UserService) ProfileLogic(id uuid.UUID) (interface{}, int) {
	user, errSelect := s.UserRepository.SelectById(id)

	if errSelect != nil || user.Username == "" {
		res := model.ResponseNotFound{
			Message: "User Not Found",
		}

		return res, http.StatusNotFound
	}

	res := model.MeProfileResponse{
		Username: user.Username,
		Name:     user.Name,
		Bio:      user.Bio,
		Web:      user.Web,
		Picture:  user.Picture,
	}

	return res, http.StatusOK
}

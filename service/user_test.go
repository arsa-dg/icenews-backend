package service

import (
	"database/sql"
	"icenews/backend/model"
	repoMock "icenews/backend/repository/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var users = []model.User{
	{
		Id:       uuid.New(),
		Username: "tester123",
		Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
		Name:     "test",
		Bio:      "test bio",
		Web:      "test web",
		Picture:  "test pict",
	},
	{
		Id:       uuid.New(),
		Username: "tester1234",
		Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
		Name:     "test1",
		Bio:      "test bio1",
		Web:      "test web1",
		Picture:  "test pict1",
	},
}

func TestService_LoginLogicOK(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester123").Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tester123",
		Password: "tester",
	})

	assert.IsType(t, model.AuthLoginResponse{}, res)
}

func TestService_LoginLogicErrorUserNotFound(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester12").Return(model.User{}, sql.ErrNoRows)

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tester12",
		Password: "tester",
	})

	assert.IsType(t, model.ResponseUnauthorized{}, res)
}

func TestService_LoginLogicErrorWrongPassword(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester123").Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tester123",
		Password: "tester123",
	})

	assert.IsType(t, model.ResponseUnauthorized{}, res)
}

func TestService_LoginLogicErrorValidation(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tes",
		Password: "t",
	})

	assert.IsType(t, model.ResponseValidationFailed{}, res)
}

func TestService_TokenLogicOK(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectById", mock.AnythingOfType("uuid.UUID")).Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.TokenLogic("0237d1b5-051d-41d5-b160-cde2d6ebf61a")

	assert.IsType(t, model.AuthLoginResponse{}, res)
}

func TestService_TokenLogicErrorUserNotFound(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectById", mock.AnythingOfType("uuid.UUID")).Return(model.User{}, nil)

	userService := NewUserService(userRepository)
	res, _ := userService.TokenLogic("9237d1b5-051d-41d5-b160-cde2d6ebf61b")

	assert.IsType(t, model.ResponseNotFound{}, res)
}

func TestService_RegisterLogicOK(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}

	req := model.RegisterRequest{
		Username: "tester12",
		Password: "tester",
		Name:     "a",
		Bio:      "a",
		Web:      "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
		Picture:  "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
	}

	userRepository.On("SelectByUsername", req.Username).Return(model.User{}, sql.ErrNoRows)
	userRepository.On("Insert", mock.AnythingOfType("model.User")).Return(nil)

	userService := NewUserService(userRepository)
	res, _ := userService.RegisterLogic(req)

	assert.IsType(t, model.ResponseOK{}, res)
}

func TestService_RegisterLogicErrorUsernameNotAvailable(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester123").Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.RegisterLogic(model.RegisterRequest{
		Username: "tester123",
		Password: "tester",
		Name:     "a",
		Bio:      "a",
		Web:      "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
		Picture:  "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
	})

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_RegisterLogicErrorValidation(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}

	userService := NewUserService(userRepository)
	res, _ := userService.RegisterLogic(model.RegisterRequest{
		Username: "tester12",
		Password: "tester",
		Name:     "a",
		Bio:      "a",
		Web:      "a",
		Picture:  "a",
	})

	assert.IsType(t, model.ResponseValidationFailed{}, res)
}

func TestService_ProfileLogicOK(t *testing.T) {
	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectById", users[0].Id).Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.ProfileLogic(users[0].Id)

	assert.IsType(t, model.MeProfileResponse{}, res)
}

func TestService_ProfileLogicErrorUserNotFound(t *testing.T) {
	uuid3 := uuid.New()

	userRepository := repoMock.UserRepositoryMock{}
	userRepository.On("SelectById", uuid3).Return(model.User{}, sql.ErrNoRows)

	userService := NewUserService(userRepository)
	res, _ := userService.ProfileLogic(uuid3)

	assert.IsType(t, model.ResponseNotFound{}, res)
}

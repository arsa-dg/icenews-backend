package service

import (
	"errors"
	"icenews/backend/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var uuid1 = uuid.New()
var uuid2 = uuid.New()

var users = []model.User{
	{
		Id:       uuid1,
		Username: "tester123",
		Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
		Name:     "test",
		Bio:      "test bio",
		Web:      "test web",
		Picture:  "test pict",
	},
	{
		Id:       uuid2,
		Username: "tester1234",
		Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
		Name:     "test1",
		Bio:      "test bio1",
		Web:      "test web1",
		Picture:  "test pict1",
	},
}

type UserRepositoryMock struct {
	mock.Mock
}

func (r UserRepositoryMock) SelectByUsername(username string) (model.User, error) {
	args := r.Called(username)

	return args.Get(0).(model.User), args.Error(1)
}

func (r UserRepositoryMock) SelectById(id uuid.UUID) (model.User, error) {
	args := r.Called(id)

	return args.Get(0).(model.User), args.Error(1)
}

func (r UserRepositoryMock) Insert(user model.User) error {
	args := r.Called(user)

	return args.Error(0)
}

func TestService_LoginLogicOK(t *testing.T) {
	userRepository := UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester123").Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tester123",
		Password: "tester",
	})

	assert.IsType(t, model.AuthLoginResponse{}, res)
}

func TestService_LoginLogicErrorUserNotFound(t *testing.T) {
	userRepository := UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester12").Return(model.User{}, errors.New("User not found"))

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tester12",
		Password: "tester",
	})

	assert.IsType(t, model.ResponseUnauthorized{}, res)
}

func TestService_LoginLogicErrorWrongPassword(t *testing.T) {
	userRepository := UserRepositoryMock{}
	userRepository.On("SelectByUsername", "tester123").Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tester123",
		Password: "tester123",
	})

	assert.IsType(t, model.ResponseUnauthorized{}, res)
}

func TestService_LoginLogicErrorValidation(t *testing.T) {
	userRepository := UserRepositoryMock{}

	userService := NewUserService(userRepository)
	res, _ := userService.LoginLogic(model.LoginRequest{
		Username: "tes",
		Password: "t",
	})

	assert.IsType(t, model.ResponseValidationFailed{}, res)
}

func TestService_RegisterLogicOK(t *testing.T) {
	userRepository := UserRepositoryMock{}

	req := model.RegisterRequest{
		Username: "tester12",
		Password: "tester",
		Name:     "a",
		Bio:      "a",
		Web:      "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
		Picture:  "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
	}

	userRepository.On("SelectByUsername", req.Username).Return(model.User{}, errors.New("User not found"))
	userRepository.On("Insert", mock.AnythingOfType("model.User")).Return(nil)

	userService := NewUserService(userRepository)
	res, _ := userService.RegisterLogic(req)

	assert.IsType(t, model.ResponseOK{}, res)
}

func TestService_RegisterLogicErrorUsernameNotAvailable(t *testing.T) {
	userRepository := UserRepositoryMock{}
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
	userRepository := UserRepositoryMock{}

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
	userRepository := UserRepositoryMock{}
	userRepository.On("SelectById", uuid1).Return(users[0], nil)

	userService := NewUserService(userRepository)
	res, _ := userService.ProfileLogic(uuid1)

	assert.IsType(t, model.MeProfileResponse{}, res)
}

func TestService_ProfileLogicErrorUserNotFound(t *testing.T) {
	uuid3 := uuid.New()

	userRepository := UserRepositoryMock{}
	userRepository.On("SelectById", uuid3).Return(model.User{}, errors.New("User not found"))

	userService := NewUserService(userRepository)
	res, _ := userService.ProfileLogic(uuid3)

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

package mock

import (
	"icenews/backend/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (s UserServiceMock) LoginLogic(request model.LoginRequest) (interface{}, int) {
	args := s.Called(request)

	return args.Get(0), args.Int(1)
}

func (s UserServiceMock) TokenLogic(id string) (interface{}, int) {
	args := s.Called(id)

	return args.Get(0), args.Int(1)
}

func (s UserServiceMock) RegisterLogic(request model.RegisterRequest) (interface{}, int) {
	args := s.Called(request)

	return args.Get(0), args.Int(1)
}

func (s UserServiceMock) ProfileLogic(id uuid.UUID) (interface{}, int) {
	args := s.Called(id)

	return args.Get(0), args.Int(1)
}

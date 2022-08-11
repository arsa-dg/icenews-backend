package mock

import (
	"icenews/backend/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

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

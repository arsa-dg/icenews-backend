package mock

import (
	"icenews/backend/model"
	"net/url"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type NewsServiceMock struct {
	mock.Mock
}

func (s NewsServiceMock) GetAllLogic(query url.Values) (interface{}, int) {
	args := s.Called(query)

	return args.Get(0), args.Int(1)
}

func (s NewsServiceMock) GetDetailLogic(id string) (interface{}, int) {
	args := s.Called(id)

	return args.Get(0), args.Int(1)
}

func (s NewsServiceMock) NewsCategoryLogic() (interface{}, int) {
	args := s.Called()

	return args.Get(0), args.Int(1)
}

func (s NewsServiceMock) AddCommentLogic(requestBody model.CommentRequest, newsId string, authorId uuid.UUID) (interface{}, int) {
	args := s.Called(requestBody, newsId, authorId)

	return args.Get(0), args.Int(1)
}

func (s NewsServiceMock) CommentListLogic(newsId string) (interface{}, int) {
	args := s.Called(newsId)

	return args.Get(0), args.Int(1)
}

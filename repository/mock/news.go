package mock

import (
	"icenews/backend/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type NewsRepositoryMock struct {
	mock.Mock
}

func (r NewsRepositoryMock) SelectAll(category string, scope string) ([]model.NewsListRaw, error) {
	args := r.Called(category, scope)

	return args.Get(0).([]model.NewsListRaw), args.Error(1)
}

func (r NewsRepositoryMock) SelectById(id string) ([]model.NewsDetailRaw, error) {
	args := r.Called(id)

	return args.Get(0).([]model.NewsDetailRaw), args.Error(1)
}

func (r NewsRepositoryMock) SelectAllCategory() ([]model.NewsCategory, error) {
	args := r.Called()

	return args.Get(0).([]model.NewsCategory), args.Error(1)
}

func (r NewsRepositoryMock) InsertComment(description, newsId string, authorId uuid.UUID) (int, error) {
	args := r.Called(description, newsId, authorId)

	return args.Int(0), args.Error(1)
}

func (r NewsRepositoryMock) SelectCommentByNewsId(newsId string) ([]model.Comment, error) {
	args := r.Called(newsId)

	return args.Get(0).([]model.Comment), args.Error(1)
}

package service

import (
	"icenews/backend/model"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var newsCategory1 = model.NewsCategory{
	Id:   1,
	Name: "Category1",
}

var newsCategory2 = model.NewsCategory{
	Id:   2,
	Name: "Category2",
}

var news11 = model.NewsListRaw{
	Id:              1,
	Title:           "News1",
	SlugUrl:         "news-1",
	CoverImage:      "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
	AdditionalImage: "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
	Nsfw:            false,
	CategoryId:      1,
	CategoryName:    "Category1",
	AuthorId:        users[0].Id,
	AuthorName:      users[0].Name,
	AuthorPicture:   users[0].Picture,
	Upvote:          0,
	Downvote:        0,
	View:            0,
	Comment:         0,
	CreatedAt:       "2019-08-24T14:15:22Z",
}

var news12 = model.NewsListRaw{
	Id:              1,
	Title:           "News1",
	SlugUrl:         "news-1",
	CoverImage:      "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
	AdditionalImage: "https://github.com/",
	Nsfw:            false,
	CategoryId:      1,
	CategoryName:    "Category1",
	AuthorId:        users[0].Id,
	AuthorName:      users[0].Name,
	AuthorPicture:   users[0].Picture,
	Upvote:          0,
	Downvote:        0,
	View:            0,
	Comment:         0,
	CreatedAt:       "2019-08-24T14:15:22Z",
}

var news1Content = "lorem ipsum"

var news11Detail = model.NewsDetailRaw{
	NewsListRaw: news11,
	Content:     news1Content,
}

var news12Detail = model.NewsDetailRaw{
	NewsListRaw: news12,
	Content:     news1Content,
}

var news1Comment1 = model.Comment{
	Id:          1,
	Description: "Comment 1",
	Commentator: model.Author{
		Id:      users[0].Id,
		Name:    users[0].Name,
		Picture: users[0].Picture,
	},
	CreatedAt: "2019-08-24T14:15:22Z",
}

var news1Comment2 = model.Comment{
	Id:          2,
	Description: "Comment 2",
	Commentator: model.Author{
		Id:      users[1].Id,
		Name:    users[1].Name,
		Picture: users[1].Picture,
	},
	CreatedAt: "2019-08-24T14:15:22Z",
}

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

func TestService_GetAllLogicOK(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectAll", "1", "top_news").Return([]model.NewsListRaw{news11, news12}, nil)

	newsService := NewNewsService(newsRepository)

	apicall, _ := url.Parse("https://icenews.com?category=1&scope=top_news")
	params := apicall.Query()

	res, _ := newsService.GetAllLogic(params)

	assert.IsType(t, model.NewsListResponse{}, res)
}

func TestService_GetAllLogicErrorCategoryNotInteger(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectAll", "cat1", "top_news").Return([]model.NewsListRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	apicall, _ := url.Parse("https://icenews.com?category=cat1&scope=top_news")
	params := apicall.Query()

	res, _ := newsService.GetAllLogic(params)

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_GetDetailLogicOK(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{news11Detail, news12Detail}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.GetDetailLogic("1")

	assert.IsType(t, model.NewsDetailResponse{}, res)
}

func TestService_GetDetailLogicErrorNewsNotFound(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectById", "2").Return([]model.NewsDetailRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.GetDetailLogic("2")

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_NewsCategoryLogicOK(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectAllCategory").Return([]model.NewsCategory{newsCategory1, newsCategory2}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.NewsCategoryLogic()

	assert.IsType(t, model.NewsCategoryResponse{}, res)
}

func TestService_AddCommentLogicOK(t *testing.T) {
	newsRepository := NewsRepositoryMock{}

	commentReq := model.CommentRequest{Description: "Bagus"}

	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{news11Detail, news12Detail}, nil)
	newsRepository.On("InsertComment", commentReq.Description, "1", users[1].Id).Return(3, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.AddCommentLogic(commentReq, "1", users[1].Id)

	assert.IsType(t, model.CommentAddResponse{}, res)
}

func TestService_AddCommentLogicErrorNewsNotFound(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.AddCommentLogic(model.CommentRequest{Description: "Bagus"}, "1", users[1].Id)

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_AddCommentLogicErrorValidation(t *testing.T) {
	newsRepository := NewsRepositoryMock{}

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.AddCommentLogic(model.CommentRequest{Description: ""}, "1", users[1].Id)

	assert.IsType(t, model.ResponseValidationFailed{}, res)
}

func TestService_CommentListLogicOK(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{news11Detail, news12Detail}, nil)
	newsRepository.On("SelectCommentByNewsId", "1").Return([]model.Comment{news1Comment1, news1Comment2}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.CommentListLogic("1")

	assert.IsType(t, model.CommentListResponse{}, res)
}

func TestService_CommentListLogicErrorNewsNotFound(t *testing.T) {
	newsRepository := NewsRepositoryMock{}
	newsRepository.On("SelectById", "2").Return([]model.NewsDetailRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.CommentListLogic("2")

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

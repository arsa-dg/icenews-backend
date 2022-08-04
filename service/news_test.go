package service

import (
	"icenews/backend/model"
	repoMock "icenews/backend/repository/mock"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestService_GetAllLogicOK(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectAll", "1", "top_news").Return([]model.NewsListRaw{news11, news12}, nil)

	newsService := NewNewsService(newsRepository)

	apicall, _ := url.Parse("https://icenews.com/news?category=1&scope=top_news")
	params := apicall.Query()

	res, _ := newsService.GetAllLogic(params)

	assert.IsType(t, model.NewsListResponse{}, res)
}

func TestService_GetAllLogicErrorCategoryNotInteger(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectAll", "cat1", "top_news").Return([]model.NewsListRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	apicall, _ := url.Parse("https://icenews.com/news?category=cat1&scope=top_news")
	params := apicall.Query()

	res, _ := newsService.GetAllLogic(params)

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_GetDetailLogicOK(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{news11Detail, news12Detail}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.GetDetailLogic("1")

	assert.IsType(t, model.NewsDetailResponse{}, res)
}

func TestService_GetDetailLogicErrorNewsNotFound(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectById", "2").Return([]model.NewsDetailRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.GetDetailLogic("2")

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_NewsCategoryLogicOK(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectAllCategory").Return([]model.NewsCategory{newsCategory1, newsCategory2}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.NewsCategoryLogic()

	assert.IsType(t, model.NewsCategoryResponse{}, res)
}

func TestService_AddCommentLogicOK(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}

	commentReq := model.CommentRequest{Description: "Bagus"}

	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{news11Detail, news12Detail}, nil)
	newsRepository.On("InsertComment", commentReq.Description, "1", users[1].Id).Return(3, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.AddCommentLogic(commentReq, "1", users[1].Id)

	assert.IsType(t, model.CommentAddResponse{}, res)
}

func TestService_AddCommentLogicErrorNewsNotFound(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.AddCommentLogic(model.CommentRequest{Description: "Bagus"}, "1", users[1].Id)

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

func TestService_AddCommentLogicErrorValidation(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.AddCommentLogic(model.CommentRequest{Description: ""}, "1", users[1].Id)

	assert.IsType(t, model.ResponseValidationFailed{}, res)
}

func TestService_CommentListLogicOK(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectById", "1").Return([]model.NewsDetailRaw{news11Detail, news12Detail}, nil)
	newsRepository.On("SelectCommentByNewsId", "1").Return([]model.Comment{news1Comment1, news1Comment2}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.CommentListLogic("1")

	assert.IsType(t, model.CommentListResponse{}, res)
}

func TestService_CommentListLogicErrorNewsNotFound(t *testing.T) {
	newsRepository := repoMock.NewsRepositoryMock{}
	newsRepository.On("SelectById", "2").Return([]model.NewsDetailRaw{}, nil)

	newsService := NewNewsService(newsRepository)

	res, _ := newsService.CommentListLogic("2")

	assert.IsType(t, model.ResponseBadRequest{}, res)
}

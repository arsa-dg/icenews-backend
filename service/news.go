package service

import (
	"icenews/backend/helper"
	"icenews/backend/model"
	"icenews/backend/repository"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type NewsServiceInterface interface {
	GetAllLogic(query url.Values) (interface{}, int)
	GetDetailLogic(id string) (interface{}, int)
	NewsCategoryLogic() (interface{}, int)
	AddCommentLogic(requestBody model.CommentRequest, newsId string, authorId uuid.UUID) (interface{}, int)
	CommentListLogic(newsId string) (interface{}, int)
}

type NewsService struct {
	Validator      *validator.Validate
	NewsRepository repository.NewsRepositoryInterface
}

func NewNewsService(r repository.NewsRepositoryInterface) NewsService {
	return NewsService{validator.New(), r}
}

func (s NewsService) GetAllLogic(query url.Values) (interface{}, int) {
	category := query.Get("category")
	_, errConvCategory := strconv.Atoi(strings.Replace(category, "", "0", 1))
	scope := query.Get("scope")

	if errConvCategory != nil {
		res := model.ResponseBadRequest{
			Message: "Category Must Be An Integer",
		}

		return res, http.StatusBadRequest
	}

	newsListRaw, err := s.NewsRepository.SelectAll(category, scope)

	if err != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	if len(newsListRaw) == 0 {
		res := model.NewsListResponse{
			Data: nil,
		}

		return res, http.StatusOK
	}

	sort.Slice(newsListRaw, func(i, j int) bool {
		return newsListRaw[i].Id < newsListRaw[j].Id
	})

	newsList := []model.NewsList{}
	var newsImage []string
	news := model.NewsList{}

	for _, newsRaw := range newsListRaw {
		if newsRaw.Id != news.Id {
			if news.Id != 0 {
				news.AdditionalImages = newsImage
				newsList = append(newsList, news)

				news = model.NewsList{}
			}
			newsImage = []string{}

			news.Id = newsRaw.Id
			news.Title = newsRaw.Title
			news.SlugUrl = newsRaw.SlugUrl
			news.CoverImage = newsRaw.CoverImage
			news.Nsfw = newsRaw.Nsfw
			news.CreatedAt = newsRaw.CreatedAt

			news.Category.Id = newsRaw.CategoryId
			news.Category.Name = newsRaw.CategoryName

			news.Author.Id = newsRaw.AuthorId
			news.Author.Name = newsRaw.AuthorName
			news.Author.Picture = newsRaw.AuthorPicture

			news.Counter.Upvote = newsRaw.Upvote
			news.Counter.Downvote = newsRaw.Downvote
			news.Counter.Comment = newsRaw.Comment
			news.Counter.View = newsRaw.View
		}

		if newsRaw.AdditionalImage != "" {
			newsImage = append(newsImage, newsRaw.AdditionalImage)
		}
	}

	news.AdditionalImages = newsImage
	newsList = append(newsList, news)

	res := model.NewsListResponse{
		Data: newsList,
	}

	return res, http.StatusOK
}

func (s NewsService) GetDetailLogic(id string) (interface{}, int) {
	newsDetailRaw, err := s.NewsRepository.SelectById(id)

	if err != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	if len(newsDetailRaw) == 0 {
		res := model.ResponseBadRequest{
			Message: "News Not Found",
		}

		return res, http.StatusNotFound
	}

	var newsImage []string
	news := model.NewsDetailResponse{}

	for idx, newsRaw := range newsDetailRaw {
		if newsRaw.AdditionalImage != "" {
			newsImage = append(newsImage, newsRaw.AdditionalImage)
		}

		if idx == 0 {
			news.Id = newsRaw.Id
			news.Title = newsRaw.Title
			news.Content = newsRaw.Content
			news.SlugUrl = newsRaw.SlugUrl
			news.CoverImage = newsRaw.CoverImage
			news.Nsfw = newsRaw.Nsfw
			news.CreatedAt = newsRaw.CreatedAt

			news.Category.Id = newsRaw.CategoryId
			news.Category.Name = newsRaw.CategoryName

			news.Author.Id = newsRaw.AuthorId
			news.Author.Name = newsRaw.AuthorName
			news.Author.Picture = newsRaw.AuthorPicture

			news.Counter.Upvote = newsRaw.Upvote
			news.Counter.Downvote = newsRaw.Downvote
			news.Counter.Comment = newsRaw.Comment
			news.Counter.View = newsRaw.View
		}
	}

	news.AdditionalImages = newsImage

	return news, http.StatusOK
}

func (s NewsService) NewsCategoryLogic() (interface{}, int) {
	newsCategory, err := s.NewsRepository.SelectAllCategory()

	if err != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := model.NewsCategoryResponse{
		Data: newsCategory,
	}

	return res, http.StatusOK
}

func (s NewsService) AddCommentLogic(requestBody model.CommentRequest, newsId string, authorId uuid.UUID) (interface{}, int) {
	errValidateRes, errValidateStatus := helper.RequestValidation(s.Validator, requestBody)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	newsDetailRaw, err := s.NewsRepository.SelectById(newsId)

	if err != nil || len(newsDetailRaw) == 0 {
		res := model.ResponseBadRequest{
			Message: "News Not Found",
		}

		return res, http.StatusNotFound
	}

	commentId, err := s.NewsRepository.InsertComment(requestBody.Description, newsId, authorId)

	if err != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := model.CommentAddResponse{
		Id: commentId,
	}

	return res, http.StatusOK
}

func (s NewsService) CommentListLogic(newsId string) (interface{}, int) {
	commentList, err := s.NewsRepository.SelectCommentByNewsId(newsId)

	if err != nil {
		res := model.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := model.CommentListResponse{
		Data: commentList,
	}

	return res, http.StatusOK
}

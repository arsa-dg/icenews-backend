package service

import (
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type NewsService struct {
	Validator      *validator.Validate
	NewsRepository repository.NewsRepository
}

func NewNewsService(DB *pgx.Conn) NewsService {
	return NewsService{validator.New(), repository.NewNewsRepository(DB)}
}

func (s NewsService) GetAllLogic(query url.Values) (interface{}, int) {
	category := query.Get("category")
	_, errConvCategory := strconv.Atoi(strings.Replace(category, "", "0", 1))
	scope := query.Get("scope")

	if errConvCategory != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Category Must Be An Integer",
		}

		return res, http.StatusBadRequest
	}

	newsListRaw, err := s.NewsRepository.SelectAll(category, scope)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	if len(newsListRaw) == 0 {
		res := interfaces.NewsListResponse{
			Data: nil,
		}

		return res, http.StatusOK
	}

	sort.Slice(newsListRaw, func(i, j int) bool {
		return newsListRaw[i].Id < newsListRaw[j].Id
	})

	newsList := []interfaces.NewsList{}
	var newsImage []string
	news := interfaces.NewsList{}

	for _, newsRaw := range newsListRaw {
		if newsRaw.Id != news.Id {
			if news.Id != 0 {
				news.AdditionalImages = newsImage
				newsList = append(newsList, news)

				news = interfaces.NewsList{}
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

	res := interfaces.NewsListResponse{
		Data: newsList,
	}

	return res, http.StatusOK
}

func (s NewsService) GetDetailLogic(id string) (interface{}, int) {
	newsDetailRaw, err := s.NewsRepository.SelectById(id)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	if len(newsDetailRaw) == 0 {
		res := interfaces.ResponseBadRequest{
			Message: "News Not Found",
		}

		return res, http.StatusNotFound
	}

	var newsImage []string
	news := interfaces.NewsDetailResponse{}

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
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := interfaces.NewsCategoryResponse{
		Data: newsCategory,
	}

	return res, http.StatusOK
}

func (s NewsService) AddCommentLogic(requestBody interfaces.CommentRequest, newsId string, authorId uuid.UUID) (interface{}, int) {
	errValidateRes, errValidateStatus := helper.RequestValidation(s.Validator, requestBody)

	if errValidateRes != nil {
		return errValidateRes, errValidateStatus
	}

	newsDetailRaw, err := s.NewsRepository.SelectById(newsId)

	if err != nil || len(newsDetailRaw) == 0 {
		res := interfaces.ResponseBadRequest{
			Message: "News Not Found",
		}

		return res, http.StatusNotFound
	}

	commentId, err := s.NewsRepository.InsertComment(requestBody.Description, newsId, authorId)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := interfaces.CommentAddResponse{
		Id: commentId,
	}

	return res, http.StatusOK
}

func (s NewsService) CommentListLogic(newsId string) (interface{}, int) {
	commentList, err := s.NewsRepository.SelectCommentByNewsId(newsId)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	res := interfaces.CommentListResponse{
		Data: commentList,
	}

	return res, http.StatusOK
}

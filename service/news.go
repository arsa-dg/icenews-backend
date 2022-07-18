package service

import (
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
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

	if news.Id == 0 {
		res := interfaces.NewsListResponse{
			Data: nil,
		}

		return res, http.StatusOK
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

	if news.Id == 0 {
		res := interfaces.ResponseBadRequest{
			Message: "News Not Found",
		}

		return res, http.StatusNotFound
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

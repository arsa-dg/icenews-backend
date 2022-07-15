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
	rows, err := s.NewsRepository.SelectById(id)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	newsImage := []string{}
	category := interfaces.NewsCategory{}
	author := interfaces.NewsAuthor{}
	counter := interfaces.NewsCounter{}
	news := interfaces.NewsDetailResponse{}

	count := 0
	for rows.Next() {
		var newImage string
		var errScan error

		if count < 1 {
			errScan = rows.Scan(
				&news.Id, &news.Title, &news.Content, &news.SlugUrl, &news.CoverImage,
				&newImage, &news.Nsfw, &category.Id, &category.Name, &author.Id,
				&author.Name, &author.Picture, &counter.Upvote, &counter.Downvote, &counter.Comment,
				&counter.View, &news.CreatedAt,
			)
		} else {
			errScan = rows.Scan(nil, nil, nil, nil, nil, &newImage, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		}

		if errScan != nil {
			res := interfaces.ResponseBadRequest{
				Message: "News Not Found",
			}

			return res, http.StatusNotFound
		}

		if newImage != "" {
			newsImage = append(newsImage, newImage)
		}

		count++
	}

	news.AdditionalImages = newsImage
	news.Category = category
	news.Author = author
	news.Counter = counter

	if news.Id == 0 {
		res := interfaces.ResponseBadRequest{
			Message: "News Not Found",
		}

		return res, http.StatusNotFound
	}

	return news, http.StatusOK
}

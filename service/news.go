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

	newsList := []interfaces.News{}
	var newsImage []string
	news := interfaces.News{}

	rows, err := s.NewsRepository.SelectAll(category, scope)

	if err != nil {
		res := interfaces.ResponseInternalServerError{
			Message: "Something Is Wrong",
		}

		return res, http.StatusInternalServerError
	}

	for rows.Next() {
		tempCategory := interfaces.NewsCategory{}
		tempAuthor := interfaces.NewsAuthor{}
		tempCounter := interfaces.NewsCounter{}
		tempNews := interfaces.News{}

		var newImage string
		var errScan error

		errScan = rows.Scan(
			&tempNews.Id, &tempNews.Title, &tempNews.SlugUrl, &tempNews.CoverImage,
			&newImage, &tempNews.Nsfw, &tempCategory.Id, &tempCategory.Name, &tempAuthor.Id,
			&tempAuthor.Name, &tempAuthor.Picture, &tempCounter.Upvote, &tempCounter.Downvote, &tempCounter.Comment,
			&tempCounter.View, &tempNews.CreatedAt,
		)

		if errScan != nil {
			res := interfaces.NewsListResponse{
				Data: nil,
			}

			return res, http.StatusOK
		}

		if news.Id != tempNews.Id {
			tempNews.Category = tempCategory
			tempNews.Author = tempAuthor
			tempNews.Counter = tempCounter

			if news.Id != 0 {
				news.AdditionalImages = newsImage
				newsList = append(newsList, news)
			}

			newsImage = []string{}
			news = tempNews
		}

		if newImage != "" {
			newsImage = append(newsImage, newImage)
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

package service

import (
	"icenews/backend/interfaces"
	"icenews/backend/repository"
	"net/http"

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

func (s NewsService) GetAllLogic() (interface{}, int) {
	news, err := s.NewsRepository.SelectAll()
	res := interfaces.NewsListResponse{
		Data: news,
	}

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return res, http.StatusOK
}

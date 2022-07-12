package service

import (
	"fmt"
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
	newsList, err := s.NewsRepository.SelectAll()

	fmt.Println("service")

	if err != nil {
		return nil, 400
	}
	fmt.Println(newsList)

	return newsList, http.StatusOK
}

package handler

import (
	"context"
	"icenews/backend/model"
	serviceMock "icenews/backend/service/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewsHandler_GetAllOK(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("GetAllLogic", mock.AnythingOfType("url.Values")).Return(model.NewsListResponse{
		Data: []model.NewsList{
			{
				Id:               1,
				Title:            "some news",
				SlugUrl:          "some-news",
				CoverImage:       "https://github.com/",
				AdditionalImages: []string{},
				Nsfw:             false,
				Category: model.NewsCategory{
					Id:   1,
					Name: "cat1",
				},
				Author: model.Author{
					Id:      uuid.New(),
					Name:    "some author",
					Picture: "https://github.com/",
				},
				Counter: model.NewsCounter{
					Upvote:   100,
					Downvote: 57,
					Comment:  12,
					View:     2579,
				},
				CreatedAt: "2019-08-24T14:15:22Z",
			},
		},
	}, http.StatusOK)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news?category=1&scope=top_news", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.GetAll(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestNewsHandler_GetAllError(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("GetAllLogic", mock.AnythingOfType("url.Values")).Return(model.ResponseBadRequest{
		Message: "Category Must Be An Integer",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news?category=football&scope=top_news", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.GetAll(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

func TestNewsHandler_GetDetailOK(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("GetDetailLogic", mock.Anything).Return(model.NewsDetailResponse{
		Id:               1,
		Title:            "some news",
		Content:          "some content",
		SlugUrl:          "some-news",
		CoverImage:       "https://github.com/",
		AdditionalImages: []string{},
		Nsfw:             false,
		Category: model.NewsCategory{
			Id:   1,
			Name: "cat1",
		},
		Author: model.Author{
			Id:      uuid.New(),
			Name:    "some author",
			Picture: "https://github.com/",
		},
		Counter: model.NewsCounter{
			Upvote:   100,
			Downvote: 57,
			Comment:  12,
			View:     2579,
		},
		CreatedAt: "2019-08-24T14:15:22Z",
	}, http.StatusOK)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news/1", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.GetDetail(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestNewsHandler_GetDetailError(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("GetDetailLogic", mock.Anything).Return(model.ResponseBadRequest{
		Message: "News Not Found",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news/98324542", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.GetDetail(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

func TestNewsHandler_NewsCategoryOK(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("NewsCategoryLogic").Return(model.NewsCategoryResponse{
		Data: []model.NewsCategory{
			{
				Id:   1,
				Name: "cat1",
			},
			{
				Id:   2,
				Name: "cat2",
			},
		},
	}, http.StatusOK)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news/category", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.NewsCategory(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestNewsHandler_AddCommentOK(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("AddCommentLogic", mock.AnythingOfType("model.CommentRequest"), mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(model.CommentAddResponse{
		Id: 17,
	}, http.StatusOK)

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"description": "berita hoax nih"}`)

	r := httptest.NewRequest(http.MethodGet, "/news/1/comment", body)

	// mock get user id from auth token (jwt)
	ctx := context.WithValue(r.Context(), "user_id", "9237d1b5-051d-41d5-b160-cde2d6ebf61b")
	r = r.WithContext(ctx)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.AddComment(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestNewsHandler_AddCommentError(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("AddCommentLogic", mock.AnythingOfType("model.CommentRequest"), mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(model.ResponseBadRequest{
		Message: "News Not Found",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"description": "berita hoax nih"}`)

	r := httptest.NewRequest(http.MethodGet, "/news/987435/comment", body)

	// mock get user id from auth token (jwt)
	ctx := context.WithValue(r.Context(), "user_id", "9237d1b5-051d-41d5-b160-cde2d6ebf61b")
	r = r.WithContext(ctx)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.AddComment(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

func TestNewsHandler_CommentListOK(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("CommentListLogic", mock.Anything).Return(model.CommentListResponse{
		Data: []model.Comment{
			{
				Id:          1,
				Description: "Bagus beritanya",
				Commentator: model.Author{
					Id:      uuid.New(),
					Name:    "some commentator",
					Picture: "https://github.com/",
				},
			},
		},
	}, http.StatusOK)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news/1/comment", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.CommentList(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestNewsHandler_CommentListError(t *testing.T) {
	newsService := serviceMock.NewsServiceMock{}
	newsService.On("CommentListLogic", mock.Anything).Return(model.ResponseBadRequest{
		Message: "News Not Found",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/news/98324542/comment", nil)

	newsHandler := NewNewsHandler(newsService)
	newsHandler.CommentList(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

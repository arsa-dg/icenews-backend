package handler

import (
	"context"
	"icenews/backend/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s UserServiceMock) ProfileLogic(id uuid.UUID) (interface{}, int) {
	args := s.Called(id)

	return args.Get(0), args.Int(1)
}

func TestAuthHandler_ProfileOK(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("ProfileLogic", mock.AnythingOfType("uuid.UUID")).Return(model.MeProfileResponse{
		Username: "tester123",
		Name:     "test name",
		Bio:      "Test bio",
		Web:      "https://github.com/",
		Picture:  "https://github.com/",
	}, http.StatusOK)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	r.Header.Set("Authorization", "Bearer somerandomtoken")

	ctx := context.WithValue(r.Context(), "user_id", "0237d1b5-051d-41d5-b160-cde2d6ebf61a")
	r = r.WithContext(ctx)

	meHandler := NewMeHandler(userService)
	meHandler.Profile(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestAuthHandler_ProfileError(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("ProfileLogic", mock.AnythingOfType("uuid.UUID")).Return(model.ResponseBadRequest{
		Message: "User Not Found",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	r.Header.Set("Authorization", "Bearer somerandomtoken")

	ctx := context.WithValue(r.Context(), "user_id", "9237d1b5-051d-41d5-b160-cde2d6ebf61b")
	r = r.WithContext(ctx)

	meHandler := NewMeHandler(userService)
	meHandler.Profile(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

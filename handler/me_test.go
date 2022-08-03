package handler

import (
	"context"
	"icenews/backend/model"
	serviceMock "icenews/backend/service/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMeHandler_ProfileOK(t *testing.T) {
	userService := serviceMock.UserServiceMock{}
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
	userService := serviceMock.UserServiceMock{}
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

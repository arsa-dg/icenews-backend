package handler

import (
	"context"
	"icenews/backend/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (s UserServiceMock) LoginLogic(request model.LoginRequest) (interface{}, int) {
	args := s.Called(request)

	return args.Get(0), args.Int(1)
}

func (s UserServiceMock) TokenLogic(id string) (interface{}, int) {
	args := s.Called(id)

	return args.Get(0), args.Int(1)
}

func (s UserServiceMock) RegisterLogic(request model.RegisterRequest) (interface{}, int) {
	args := s.Called(request)

	return args.Get(0), args.Int(1)
}

func TestAuthHandler_LoginOK(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("LoginLogic", model.LoginRequest{
		Username: "tester123",
		Password: "tester",
	}).Return(model.AuthLoginResponse{
		Token:      "somerandomtoken",
		Scheme:     "Bearer",
		Expires_at: "2019-08-24T14:15:22Z",
	}, http.StatusOK)

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"username": "tester123", "password": "tester"}`)

	r := httptest.NewRequest(http.MethodPost, "/auth/login", body)

	authHandler := NewAuthHandler(userService)
	authHandler.Login(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestAuthHandler_LoginError(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("LoginLogic", model.LoginRequest{
		Username: "tester123",
		Password: "testttt",
	}).Return(model.ResponseUnauthorized{
		Message: "Wrong Password",
	}, http.StatusUnauthorized)

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"username": "tester123", "password": "testttt"}`)

	r := httptest.NewRequest(http.MethodPost, "/auth/login", body)

	authHandler := NewAuthHandler(userService)
	authHandler.Login(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusUnauthorized, resStatusCode)
}

func TestAuthHandler_TokenOK(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("TokenLogic", "someuserid").Return(model.AuthLoginResponse{
		Token:      "someNEWrandomtoken",
		Scheme:     "Bearer",
		Expires_at: "2029-08-24T14:15:22Z",
	}, http.StatusOK)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	r.Header.Set("Authorization", "Bearer somerandomtoken")

	// mock get user id from auth token (jwt)
	ctx := context.WithValue(r.Context(), "user_id", "someuserid")
	r = r.WithContext(ctx)

	authHandler := NewAuthHandler(userService)
	authHandler.Token(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestAuthHandler_TokenError(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("TokenLogic", "someuserid").Return(model.ResponseBadRequest{
		Message: "User Not Found",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	r.Header.Set("Authorization", "Bearer somerandomtoken")

	ctx := context.WithValue(r.Context(), "user_id", "someuserid")
	r = r.WithContext(ctx)

	authHandler := NewAuthHandler(userService)
	authHandler.Token(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

func TestAuthHandler_RegisterOK(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("RegisterLogic", model.RegisterRequest{
		Username: "tester12",
		Password: "tester",
		Name:     "a",
		Bio:      "a",
		Web:      "https://github.com/",
		Picture:  "https://github.com/",
	}).Return(model.ResponseOK{
		Message: "Register Success",
	}, http.StatusOK)

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"username": "tester12", "password": "tester", "name": "a", "bio": "a", "web": "https://github.com/", "picture": "https://github.com/"}`)

	r := httptest.NewRequest(http.MethodPost, "/auth/register", body)

	authHandler := NewAuthHandler(userService)
	authHandler.Register(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusOK, resStatusCode)
}

func TestAuthHandler_RegisterError(t *testing.T) {
	userService := UserServiceMock{}
	userService.On("RegisterLogic", model.RegisterRequest{
		Username: "tester123",
		Password: "tester",
		Name:     "a",
		Bio:      "a",
		Web:      "https://github.com/",
		Picture:  "https://github.com/",
	}).Return(model.ResponseBadRequest{
		Message: "Username Is Not Available",
	}, http.StatusBadRequest)

	w := httptest.NewRecorder()

	body := strings.NewReader(`{"username": "tester123", "password": "tester", "name": "a", "bio": "a", "web": "https://github.com/", "picture": "https://github.com/"}`)

	r := httptest.NewRequest(http.MethodPost, "/auth/register", body)

	authHandler := NewAuthHandler(userService)
	authHandler.Register(w, r)

	resStatusCode := w.Result().StatusCode

	assert.Equal(t, http.StatusBadRequest, resStatusCode)
}

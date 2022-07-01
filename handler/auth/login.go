package auth

import (
	"encoding/json"
	"icenews/backend/helper"
	"icenews/backend/repository"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginField struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type responseOK struct {
	Token      string `json:"token"`
	Scheme     string `json:"scheme"`
	Expires_at string `json:"expires_at"`
}

func (AH AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// to do
	// validate hash salt?
	// validate username/password empty?
	// response 400, 401, 422

	var field LoginField
	json.NewDecoder(r.Body).Decode(&field)

	userRepository := repository.NewUserRepository(AH.DB)

	user := userRepository.SelectByUsername(field.Username)

	if user.Password == field.Password {
		token, expiresAt := createJWT(user.Id)

		res := responseOK{
			Token:      token,
			Scheme:     "Bearer",
			Expires_at: expiresAt,
		}

		helper.ResponseOK(w, res)
	}
}

func createJWT(id string) (string, string) {
	expiresAt := time.Now().UTC().Add(time.Hour * 2) // 2 hours
	jwtExp := expiresAt.Unix()

	expiresAtString := expiresAt.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     jwtExp,
	})

	tokenString, _ := token.SignedString([]byte("SECRET")) // sementara hardcoded

	return tokenString, expiresAtString
}

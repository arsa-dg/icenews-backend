package middleware

import (
	"context"
	"fmt"
	"icenews/backend/helper"
	"icenews/backend/model"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string
		secretKey := os.Getenv("SECRET_KEY")

		auth := r.Header.Get("Authorization")

		if auth == "" {
			res := model.ResponseUnauthorized{
				Message: "Authorization is missing",
			}

			helper.ResponseError(w, http.StatusUnauthorized, res)
		} else {
			authSplit := strings.Split(auth, " ")
			tokenString := authSplit[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(secretKey), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userId = claims["user_id"].(string)

				ctx := context.WithValue(r.Context(), "user_id", userId)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				res := model.ResponseUnauthorized{
					Message: err.Error(),
				}

				helper.ResponseError(w, http.StatusUnauthorized, res)
			}
		}
	})
}

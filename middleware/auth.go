package middleware

import (
	"context"
	"fmt"
	"icenews/backend/helper"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type responseUnauthorized struct {
	Message string `json:"message"`
}

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string

		auth := r.Header.Get("Authorization")

		if auth == "" {
			res := responseUnauthorized{
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

				return []byte("SECRET"), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userId = claims["user_id"].(string)
			} else {
				res := responseUnauthorized{
					Message: err.Error(),
				}

				helper.ResponseError(w, http.StatusUnauthorized, res)
			}

			ctx := context.WithValue(r.Context(), "user_id", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

func CreateJWT(id string) (string, string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	expiresAt := time.Now().UTC().Add(time.Hour * 2) // 2 hours
	jwtExp := expiresAt.Unix()

	expiresAtString := expiresAt.Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     jwtExp,
	})

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Error().Err(err).Msg("Error while creating JWT token")
	}

	return tokenString, expiresAtString, err
}

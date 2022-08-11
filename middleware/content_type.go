package middleware

import (
	"errors"
	"icenews/backend/helper"
	"icenews/backend/model"
	"mime"
	"net/http"

	"github.com/rs/zerolog/log"
)

func TypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType == "" {
			log.Error().Err(errors.New("content type is missing")).Msg("Error content type is missing")

			res := model.ResponseBadRequest{
				Message: "Content-Type is missing",
			}

			helper.ResponseError(w, http.StatusBadRequest, res)
		} else {
			mediaType, _, err := mime.ParseMediaType(contentType)

			if err != nil {
				log.Error().Err(err).Msg("Error content type")

				res := model.ResponseBadRequest{
					Message: "Content-Type is malformed",
				}

				helper.ResponseError(w, http.StatusBadRequest, res)

				return
			}

			if mediaType != "application/json" {
				log.Error().Err(errors.New("content type is incorrect")).Msg("Error content type is not application/json")

				res := model.ResponseBadRequest{
					Message: "Content-Type is in wrong format",
				}

				helper.ResponseError(w, http.StatusBadRequest, res)

				return
			}

			next.ServeHTTP(w, r)
		}
	})
}

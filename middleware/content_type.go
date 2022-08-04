package middleware

import (
	"icenews/backend/helper"
	"icenews/backend/model"
	"mime"
	"net/http"
)

func TypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType == "" {
			res := model.ResponseBadRequest{
				Message: "Content-Type is missing",
			}

			helper.ResponseError(w, http.StatusBadRequest, res)
		} else {
			mediaType, _, err := mime.ParseMediaType(contentType)

			if err != nil {
				res := model.ResponseBadRequest{
					Message: "Content-Type is malformed",
				}

				helper.ResponseError(w, http.StatusBadRequest, res)

				return
			}

			if mediaType != "application/json" {
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

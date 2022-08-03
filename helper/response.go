package helper

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func ResponseOK(w http.ResponseWriter, data interface{}) {
	response(w, http.StatusOK, data)
}

func ResponseError(w http.ResponseWriter, httpCode int, data interface{}) {
	response(w, httpCode, data)
}

func response(w http.ResponseWriter, httpCode int, data interface{}) {
	log.Info().
		Str("status", http.StatusText(httpCode)).
		Msg("response")

	response, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(response)
}

package auth

import (
	"icenews/backend/helper"
	"icenews/backend/interfaces"
	"net/http"
)

func (AH AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	token, expiresAt, err := helper.CreateJWT(userId)

	// bad request (400)
	if err != nil {
		res := interfaces.ResponseBadRequest{
			Message: "Something Is Wrong",
		}

		helper.ResponseError(w, http.StatusBadRequest, res)
	}

	res := interfaces.AuthResponseOK{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	helper.ResponseOK(w, res)
}

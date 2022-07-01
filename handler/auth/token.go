package auth

import (
	"icenews/backend/helper"
	"net/http"
)

func (AH AuthHandler) Token(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	token, expiresAt := helper.CreateJWT(userId)

	res := responseOK{
		Token:      token,
		Scheme:     "Bearer",
		Expires_at: expiresAt,
	}

	helper.ResponseOK(w, res)
}

package auth

import "github.com/jackc/pgx/v4"

type AuthHandler struct {
	DB *pgx.Conn
}

func NewAuthHandler(DB *pgx.Conn) AuthHandler {
	return AuthHandler{DB}
}

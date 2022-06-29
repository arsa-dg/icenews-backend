package handler

import "github.com/jackc/pgx/v4"

type handler struct {
	DB *pgx.Conn
}

func New(DB *pgx.Conn) handler {
	return handler{DB}
}

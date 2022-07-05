package repository

import (
	"context"
	"icenews/backend/interfaces"

	"github.com/jackc/pgx/v4"
)

type UserRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(DB *pgx.Conn) UserRepository {
	return UserRepository{DB}
}

func (Repository UserRepository) SelectByUsername(username string) interfaces.User {
	user := interfaces.User{}

	Repository.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Bio,
		&user.Web,
		&user.Picture,
	)

	return user
}

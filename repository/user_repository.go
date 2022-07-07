package repository

import (
	"context"
	"icenews/backend/interfaces"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type UserRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(DB *pgx.Conn) UserRepository {
	return UserRepository{DB}
}

func (Repository UserRepository) SelectByUsername(username string) (interfaces.User, error) {
	user := interfaces.User{}

	err := Repository.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Bio,
		&user.Web,
		&user.Picture,
	)

	return user, err
}

func (Repository UserRepository) Insert(username string, password string, name string, bio string, web string, picture string) error {
	id := uuid.New()

	_, err := Repository.DB.Exec(context.Background(), `INSERT INTO
		users(
			id, username, password, name, bio, web, picture
		) 
		values(
			$1, $2, $3, $4, $5, $6, $7
		)
	`, id, username, password, name, bio, web, picture)

	return err
}

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

func (r UserRepository) SelectByUsername(username string) (interfaces.User, error) {
	user := interfaces.User{}

	err := r.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE username=$1", username).Scan(
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

func (r UserRepository) SelectById(id uuid.UUID) (interfaces.User, error) {
	user := interfaces.User{}

	err := r.DB.QueryRow(context.Background(), "SELECT * FROM users WHERE id=$1", id).Scan(
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

func (r UserRepository) Insert(user interfaces.User) error {
	_, err := r.DB.Exec(context.Background(), `INSERT INTO
		users(
			id, username, password, name, bio, web, picture
		) 
		values(
			$1, $2, $3, $4, $5, $6, $7
		)
	`, user.Id, user.Username, user.Password, user.Name, user.Bio, user.Web, user.Picture)

	return err
}

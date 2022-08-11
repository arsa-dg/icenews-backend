package repository

import (
	"context"
	"database/sql"
	"icenews/backend/model"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserRepositoryInterface interface {
	SelectByUsername(username string) (model.User, error)
	SelectById(id uuid.UUID) (model.User, error)
	Insert(user model.User) error
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) UserRepository {
	return UserRepository{DB}
}

func (r UserRepository) SelectByUsername(username string) (model.User, error) {
	user := model.User{}

	err := r.DB.QueryRowContext(context.Background(), "SELECT * FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Bio,
		&user.Web,
		&user.Picture,
	)

	if err != nil {
		log.Error().Err(err).Msg("Error select user by username")
	}

	return user, err
}

func (r UserRepository) SelectById(id uuid.UUID) (model.User, error) {
	user := model.User{}

	err := r.DB.QueryRowContext(context.Background(), "SELECT * FROM users WHERE id=$1", id).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Name,
		&user.Bio,
		&user.Web,
		&user.Picture,
	)

	if err != nil {
		log.Error().Err(err).Msg("Error select user by id")
	}

	return user, err
}

func (r UserRepository) Insert(user model.User) error {
	_, err := r.DB.ExecContext(context.Background(), `INSERT INTO
		users(
			id, username, password, name, bio, web, picture
		) 
		values(
			$1, $2, $3, $4, $5, $6, $7
		)
	`, user.Id, user.Username, user.Password, user.Name, user.Bio, user.Web, user.Picture)

	if err != nil {
		log.Error().Err(err).Msg("Error insert user")
	}

	return err
}

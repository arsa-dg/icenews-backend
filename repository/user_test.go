package repository

import (
	"database/sql"
	"errors"
	"icenews/backend/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var user1 = model.User{
	Id:       uuid.New(),
	Username: "tester123",
	Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
	Name:     "test",
	Bio:      "test bio",
	Web:      "test web",
	Picture:  "test pict",
}

func TestUserRepository_SelectByUsernameOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"id", "username", "password", "name", "bio", "web", "picture"})
	rows.AddRow(user1.Id, user1.Username, user1.Password, user1.Name, user1.Bio, user1.Web, user1.Picture)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username=$1")).WithArgs(user1.Username).WillReturnRows(rows)

	userRepository := NewUserRepository(DB)
	_, errSelect := userRepository.SelectByUsername(user1.Username)

	assert.Nil(t, errSelect)
}

func TestUserRepository_SelectByUsernameErrorUserNotFound(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"id", "username", "password", "name", "bio", "web", "picture"})
	rows.AddRow(user1.Id, user1.Username, user1.Password, user1.Name, user1.Bio, user1.Web, user1.Picture)

	uname := "usernamebohong"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username=$1")).WithArgs(uname).WillReturnError(sql.ErrNoRows)

	userRepository := NewUserRepository(DB)
	_, errSelect := userRepository.SelectByUsername(uname)

	assert.Error(t, errSelect)
}

func TestUserRepository_SelectByIdOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"id", "username", "password", "name", "bio", "web", "picture"})
	rows.AddRow(user1.Id, user1.Username, user1.Password, user1.Name, user1.Bio, user1.Web, user1.Picture)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id=$1")).WithArgs(user1.Id).WillReturnRows(rows)

	userRepository := NewUserRepository(DB)
	_, errSelect := userRepository.SelectById(user1.Id)

	assert.Nil(t, errSelect)
}

func TestUserRepository_SelectByIdErrorUserNotFound(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"id", "username", "password", "name", "bio", "web", "picture"})
	rows.AddRow(user1.Id, user1.Username, user1.Password, user1.Name, user1.Bio, user1.Web, user1.Picture)

	id := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username=$1")).WithArgs(id).WillReturnError(sql.ErrNoRows)

	userRepository := NewUserRepository(DB)
	_, errSelect := userRepository.SelectById(id)

	assert.Error(t, errSelect)
}

func TestUserRepository_InsertOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	newUser := model.User{
		Id:       uuid.New(),
		Username: "tester1234",
		Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
		Name:     "test2",
		Bio:      "test2 bio",
		Web:      "test2 web",
		Picture:  "test2 pict",
	}

	rows := mock.NewRows([]string{"id", "username", "password", "name", "bio", "web", "picture"})
	rows.AddRow(user1.Id, user1.Username, user1.Password, user1.Name, user1.Bio, user1.Web, user1.Picture)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO
		users(
			id, username, password, name, bio, web, picture
		) 
		values(
			$1, $2, $3, $4, $5, $6, $7
		)
	`)).WithArgs(newUser.Id, newUser.Username, newUser.Password, newUser.Name, newUser.Bio, newUser.Web, newUser.Picture).WillReturnResult(sqlmock.NewResult(1, 1))

	userRepository := NewUserRepository(DB)
	errInsert := userRepository.Insert(newUser)

	assert.Nil(t, errInsert)
}

func TestUserRepository_InsertErrorUniqueConstraintViolated(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	newUser := model.User{
		Id:       uuid.New(),
		Username: "tester123",
		Password: "$2a$10$NEQETyrW4pBS1e/dX1DSAOQoZaD/x./sNm5PJQuz34BeGt6Y5b3Zm",
		Name:     "test2",
		Bio:      "test2 bio",
		Web:      "test2 web",
		Picture:  "test2 pict",
	}

	rows := mock.NewRows([]string{"id", "username", "password", "name", "bio", "web", "picture"})
	rows.AddRow(user1.Id, user1.Username, user1.Password, user1.Name, user1.Bio, user1.Web, user1.Picture)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO
		users(
			id, username, password, name, bio, web, picture
		) 
		values(
			$1, $2, $3, $4, $5, $6, $7
		)
	`)).WithArgs(newUser.Id, newUser.Username, newUser.Password, newUser.Name, newUser.Bio, newUser.Web, newUser.Picture).
		WillReturnError(errors.New(`duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)`))

	userRepository := NewUserRepository(DB)
	errInsert := userRepository.Insert(newUser)

	assert.Error(t, errInsert)
}

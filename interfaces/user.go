package interfaces

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Username string
	Password string
	Name     string
	Bio      string
	Web      string
	Picture  string
}

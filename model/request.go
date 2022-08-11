package model

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=100"`
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio" validate:"required"`
	Web      string `json:"web" validate:"required,uri"`
	Picture  string `json:"picture" validate:"required,uri"`
}

type CommentRequest struct {
	Description string `json:"description" validate:"required,min=1,max=255"`
}

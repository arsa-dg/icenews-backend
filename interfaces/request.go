package interfaces

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=4"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=4"`
	Name     string `json:"name" validate:"required"`
	Bio      string `json:"bio" validate:"required"`
	Web      string `json:"web" validate:"required,uri"`
	Picture  string `json:"picture" validate:"required,uri"`
}

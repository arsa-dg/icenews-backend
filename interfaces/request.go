package interfaces

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
	Picture  string `json:"picture"`
}

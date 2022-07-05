package interfaces

type User struct {
	Id       string
	Username string
	Password string
	Name     string
	Bio      string
	Web      string
	Picture  string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

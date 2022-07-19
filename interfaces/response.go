package interfaces

import "github.com/google/uuid"

type AuthLoginResponse struct {
	Token      string `json:"token"`
	Scheme     string `json:"scheme"`
	Expires_at string `json:"expires_at"`
}

type MeProfileResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
	Picture  string `json:"picture"`
}

type NewsCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NewsAuthor struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Picture string    `json:"picture"`
}

type NewsCounter struct {
	Upvote   int `json:"upvote"`
	Downvote int `json:"downvote"`
	Comment  int `json:"comment"`
	View     int `json:"view"`
}

type NewsList struct {
	Id               int          `json:"id"`
	Title            string       `json:"title"`
	SlugUrl          string       `json:"slug_url"`
	CoverImage       string       `json:"cover_image"`
	AdditionalImages []string     `json:"additional_images"`
	Nsfw             bool         `json:"nsfw"`
	Category         NewsCategory `json:"category"`
	Author           NewsAuthor   `json:"author"`
	Counter          NewsCounter  `json:"counter"`
	CreatedAt        string       `json:"created_at"`
}

type NewsDetailResponse struct {
	Id               int          `json:"id"`
	Title            string       `json:"title"`
	Content          string       `json:"content"`
	SlugUrl          string       `json:"slug_url"`
	CoverImage       string       `json:"cover_image"`
	AdditionalImages []string     `json:"additional_images"`
	Nsfw             bool         `json:"nsfw"`
	Category         NewsCategory `json:"category"`
	Author           NewsAuthor   `json:"author"`
	Counter          NewsCounter  `json:"counter"`
	CreatedAt        string       `json:"created_at"`
}

type NewsListResponse struct {
	Data []NewsList `json:"data"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseBadRequest struct {
	Message string `json:"message"`
}

type ResponseUnauthorized struct {
	Message string `json:"message"`
}

type FieldError struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

type ResponseValidationFailed struct {
	Message string       `json:"message"`
	Field   []FieldError `json:"field"`
}

type ResponseInternalServerError struct {
	Message string `json:"message"`
}

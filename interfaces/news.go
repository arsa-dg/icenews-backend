package interfaces

import "github.com/google/uuid"

type NewsCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Author struct {
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
	Author           Author       `json:"author"`
	Counter          NewsCounter  `json:"counter"`
	CreatedAt        string       `json:"created_at"`
}

type NewsListRaw struct {
	Id              int
	Title           string
	SlugUrl         string
	CoverImage      string
	AdditionalImage string
	Nsfw            bool
	CategoryId      int
	CategoryName    string
	AuthorId        uuid.UUID
	AuthorName      string
	AuthorPicture   string
	Upvote          int
	Downvote        int
	Comment         int
	View            int
	CreatedAt       string
}

type NewsDetailRaw struct {
	NewsListRaw
	Content string
}

type Comment struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Commentator Author `json:"commentator"`
	CreatedAt   string `json:"created_at"`
}

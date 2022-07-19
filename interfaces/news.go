package interfaces

import "github.com/google/uuid"

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

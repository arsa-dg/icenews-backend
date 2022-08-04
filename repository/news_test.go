package repository

import (
	"database/sql"
	"icenews/backend/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var newsCategories = []model.NewsCategory{
	{
		Id:   1,
		Name: "Cat1",
	},
	{
		Id:   2,
		Name: "Cat2",
	},
}

var newsDetail1 = model.NewsDetailRaw{
	NewsListRaw: model.NewsListRaw{
		Id:              1,
		Title:           "News1",
		SlugUrl:         "news-1",
		CoverImage:      "https://github.com/JuanCrg90/Clean-Code-Notes#chapter9",
		AdditionalImage: "https://github.com/",
		Nsfw:            false,
		CategoryId:      newsCategories[0].Id,
		CategoryName:    newsCategories[0].Name,
		AuthorId:        user1.Id,
		AuthorName:      user1.Name,
		AuthorPicture:   user1.Picture,
		Upvote:          0,
		Downvote:        0,
		View:            0,
		Comment:         0,
		CreatedAt:       "2019-08-24T14:15:22Z",
	},
	Content: "berita hari ini adalah",
}

var news1Comments = []model.Comment{
	{
		Id:          1,
		Description: "bagus",
		Commentator: model.Author{
			Id:      user1.Id,
			Name:    user1.Name,
			Picture: user1.Picture,
		},
	},
	{
		Id:          2,
		Description: "good",
		Commentator: model.Author{
			Id:      user1.Id,
			Name:    user1.Name,
			Picture: user1.Picture,
		},
	},
}

func TestNewsRepository_SelectAllOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"news.id", "news.title", "news.slug_url", "news.cover_image", "coalesce", "news.nsfw", "categories.id", "categories.name", "users.id", "users.name",
		"users.picture", "news.upvote", "news.downvote", "news.comment", "news.view", "TO_CHAR"})
	rows.AddRow(newsDetail1.Id, newsDetail1.Title, newsDetail1.SlugUrl, newsDetail1.CoverImage, newsDetail1.AdditionalImage, newsDetail1.Nsfw, newsDetail1.CategoryId,
		newsDetail1.CategoryName, newsDetail1.AuthorId, newsDetail1.AuthorName, newsDetail1.AuthorPicture, newsDetail1.Upvote, newsDetail1.Downvote, newsDetail1.Comment, newsDetail1.View, newsDetail1.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT
		news.id, news.title, news.slug_url, news.cover_image, 
		COALESCE(news_images.image, ''), news.nsfw, categories.id, categories.name, 
		users.id, users.name, users.picture, news.upvote, news.downvote, 
		news.comment, news.view, TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		LEFT JOIN news_images ON news.id = news_images.news_id
		WHERE
		news.category_id = $1 AND
		news.scope = $2
	`)).WithArgs("1", "top_news").WillReturnRows(rows)

	newsRepository := NewNewsRepository(DB)
	result, _ := newsRepository.SelectAll("1", "top_news")

	assert.Equal(t, []model.NewsListRaw{newsDetail1.NewsListRaw}, result)
}

func TestNewsRepository_SelectByIdOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"news.id", "news.title", "news.content", "news.slug_url", "news.cover_image", "coalesce", "news.nsfw", "categories.id", "categories.name", "users.id", "users.name",
		"users.picture", "news.upvote", "news.downvote", "news.comment", "news.view", "TO_CHAR"})
	rows.AddRow(newsDetail1.Id, newsDetail1.Title, newsDetail1.Content, newsDetail1.SlugUrl, newsDetail1.CoverImage, newsDetail1.AdditionalImage, newsDetail1.Nsfw, newsDetail1.CategoryId,
		newsDetail1.CategoryName, newsDetail1.AuthorId, newsDetail1.AuthorName, newsDetail1.AuthorPicture, newsDetail1.Upvote, newsDetail1.Downvote, newsDetail1.Comment, newsDetail1.View, newsDetail1.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT
		news.id, news.title, news.content, news.slug_url, news.cover_image, 
		COALESCE(news_images.image, ''), news.nsfw, categories.id, categories.name, 
		users.id, users.name, users.picture, news.upvote, news.downvote, 
		news.comment, news.view, TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		LEFT JOIN news_images ON news.id = news_images.news_id
		WHERE news.id = $1;
	`)).WithArgs("1").WillReturnRows(rows)

	newsRepository := NewNewsRepository(DB)
	result, _ := newsRepository.SelectById("1")

	assert.Equal(t, []model.NewsDetailRaw{newsDetail1}, result)
}

func TestNewsRepository_SelectByIdErrorNewsNotFound(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT
		news.id, news.title, news.content, news.slug_url, news.cover_image, 
		COALESCE(news_images.image, ''), news.nsfw, categories.id, categories.name, 
		users.id, users.name, users.picture, news.upvote, news.downvote, 
		news.comment, news.view, TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		LEFT JOIN news_images ON news.id = news_images.news_id
		WHERE news.id = $1;
	`)).WithArgs("2").WillReturnError(sql.ErrNoRows)

	newsRepository := NewNewsRepository(DB)
	_, errSelect := newsRepository.SelectById("2")

	assert.Error(t, errSelect)
}

func TestUserRepository_InsertCommentOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"id"})
	rows.AddRow(3)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO comments(description, author_id, news_id) values($1, $2, $3) RETURNING id")).
		WithArgs("oke", "1", user1.Id).WillReturnRows(rows)

	newsRepository := NewNewsRepository(DB)
	result, _ := newsRepository.InsertComment("oke", "1", user1.Id)

	assert.NotNil(t, result)
}

func TestUserRepository_SelectAllCategoryOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"id", "name"})
	rows.AddRow(newsCategories[0].Id, newsCategories[0].Name)
	rows.AddRow(newsCategories[1].Id, newsCategories[1].Name)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM categories")).WillReturnRows(rows)

	newsRepository := NewNewsRepository(DB)
	result, _ := newsRepository.SelectAllCategory()

	assert.Equal(t, newsCategories, result)
}

func TestUserRepository_SelectCommentByNewsIdOK(t *testing.T) {
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Error while mocking database")
	}

	defer DB.Close()

	rows := mock.NewRows([]string{"comments.id", "comments.description", "users.id", "users.name", "users.picture", "to.char"})
	rows.AddRow(news1Comments[0].Id, news1Comments[0].Description, news1Comments[0].Commentator.Id, news1Comments[0].Commentator.Name, news1Comments[0].Commentator.Picture, news1Comments[0].CreatedAt)
	rows.AddRow(news1Comments[1].Id, news1Comments[1].Description, news1Comments[1].Commentator.Id, news1Comments[1].Commentator.Name, news1Comments[1].Commentator.Picture, news1Comments[1].CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT
		comments.id, comments.description, users.id, users.name, users.picture, 
		TO_CHAR(comments.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM comments 
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.news_id = $1;
	`)).WithArgs("1").WillReturnRows(rows)

	newsRepository := NewNewsRepository(DB)
	result, _ := newsRepository.SelectCommentByNewsId("1")

	assert.Equal(t, news1Comments, result)
}

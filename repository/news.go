package repository

import (
	"context"
	"fmt"
	"icenews/backend/interfaces"

	"github.com/jackc/pgx/v4"
)

type NewsRepository struct {
	DB *pgx.Conn
}

func NewNewsRepository(DB *pgx.Conn) NewsRepository {
	return NewsRepository{DB}
}

func (r NewsRepository) SelectAll() (interfaces.NewsList, error) {
	fmt.Println("repository")

	var newsImage []string
	category := interfaces.NewsCategory{}
	author := interfaces.NewsAuthor{}
	counter := interfaces.NewsCounter{}
	newsList := interfaces.NewsList{}

	rows, err := r.DB.Query(context.Background(), `SELECT 
		news.id, news.title, news.slug_url, news.cover_image, news_images.image, news.nsfw, 
		categories.id, categories.name, users.id, users.name, users.picture, 
		news.upvote, news.downvote, news.comment, news.view, 
		TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		INNER JOIN news_images ON news.id = news_images.news_id;
	`)

	count := 0
	for rows.Next() {
		var newImage string
		var errScan error
		fmt.Println(count)

		if count < 1 {
			errScan = rows.Scan(
				&newsList.Id, &newsList.Title, &newsList.SlugUrl, &newsList.CoverImage,
				&newImage, &newsList.Nsfw, &category.Id, &category.Name, &author.Id,
				&author.Name, &author.Picture, &counter.Upvote, &counter.Downvote, &counter.Comment,
				&counter.View, &newsList.CreatedAt,
			)
		} else {
			errScan = rows.Scan(nil, nil, nil, nil, &newImage, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		}

		if errScan != nil {
			fmt.Println(errScan)
			return newsList, errScan
		}

		newsImage = append(newsImage, newImage)

		count++
	}

	fmt.Println(newsList)

	newsList.AdditionalImages = newsImage
	newsList.Category = category
	newsList.Author = author
	newsList.Counter = counter

	return newsList, err
}

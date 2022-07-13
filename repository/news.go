package repository

import (
	"context"
	"icenews/backend/interfaces"

	"github.com/jackc/pgx/v4"
)

type NewsRepository struct {
	DB *pgx.Conn
}

func NewNewsRepository(DB *pgx.Conn) NewsRepository {
	return NewsRepository{DB}
}

func (r NewsRepository) SelectAll() ([]interfaces.News, error) {
	newsList := []interfaces.News{}

	var newsImage []string
	news := interfaces.News{}

	rows, err := r.DB.Query(context.Background(), `SELECT 
		news.id, news.title, news.slug_url, news.cover_image, COALESCE(news_images.image, ''), 
		news.nsfw, categories.id, categories.name, users.id, users.name, users.picture, 
		news.upvote, news.downvote, news.comment, news.view, 
		TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		LEFT JOIN news_images ON news.id = news_images.news_id;
	`)

	for rows.Next() {
		tempCategory := interfaces.NewsCategory{}
		tempAuthor := interfaces.NewsAuthor{}
		tempCounter := interfaces.NewsCounter{}
		tempNews := interfaces.News{}

		var newImage string
		var errScan error

		errScan = rows.Scan(
			&tempNews.Id, &tempNews.Title, &tempNews.SlugUrl, &tempNews.CoverImage,
			&newImage, &tempNews.Nsfw, &tempCategory.Id, &tempCategory.Name, &tempAuthor.Id,
			&tempAuthor.Name, &tempAuthor.Picture, &tempCounter.Upvote, &tempCounter.Downvote, &tempCounter.Comment,
			&tempCounter.View, &tempNews.CreatedAt,
		)

		if errScan != nil {
			return newsList, errScan
		}

		if news.Id != tempNews.Id {
			tempNews.Category = tempCategory
			tempNews.Author = tempAuthor
			tempNews.Counter = tempCounter

			if news.Id != 0 {
				news.AdditionalImages = newsImage
				newsList = append(newsList, news)
			}

			newsImage = []string{}
			news = tempNews
		}

		if newImage != "" {
			newsImage = append(newsImage, newImage)
		}
	}

	news.AdditionalImages = newsImage
	newsList = append(newsList, news)

	return newsList, err
}

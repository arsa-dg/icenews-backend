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

func (r NewsRepository) SelectAll(category string, scope string) ([]interfaces.NewsListRaw, error) {
	var rows pgx.Rows
	var err error
	var newsListRaw []interfaces.NewsListRaw

	query := `SELECT
		news.id, news.title, news.slug_url, news.cover_image, COALESCE(news_images.image, ''),
		news.nsfw, categories.id, categories.name, users.id, users.name, users.picture,
		news.upvote, news.downvote, news.comment, news.view,
		TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		LEFT JOIN news_images ON news.id = news_images.news_id`

	if category != "" && scope != "" {
		rows, err = r.DB.Query(context.Background(), query+` WHERE
			news.category_id = $1 AND
			news.scope = $2;
		`, category, scope)
	} else if category != "" {
		rows, err = r.DB.Query(context.Background(), query+` WHERE
			news.category_id = $1;
		`, category)
	} else if scope != "" {
		rows, err = r.DB.Query(context.Background(), query+` WHERE
			news.scope = $1;
		`, scope)
	} else {
		rows, err = r.DB.Query(context.Background(), query)
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		news := interfaces.NewsListRaw{}

		errScan := rows.Scan(
			&news.Id, &news.Title, &news.SlugUrl, &news.CoverImage,
			&news.AdditionalImage, &news.Nsfw, &news.CategoryId,
			&news.CategoryName, &news.AuthorId, &news.AuthorName,
			&news.AuthorPicture, &news.Upvote, &news.Downvote,
			&news.Comment, &news.View, &news.CreatedAt,
		)

		if errScan != nil {
			return nil, errScan
		}

		newsListRaw = append(newsListRaw, news)
	}

	return newsListRaw, err
}

func (r NewsRepository) SelectById(id string) (interfaces.NewsDetailResponse, error) {
	newsImage := []string{}
	category := interfaces.NewsCategory{}
	author := interfaces.NewsAuthor{}
	counter := interfaces.NewsCounter{}
	news := interfaces.NewsDetailResponse{}

	rows, err := r.DB.Query(context.Background(), `SELECT
		news.id, news.title, news.content, news.slug_url, news.cover_image, 
		COALESCE(news_images.image, ''), news.nsfw, categories.id, categories.name, 
		users.id, users.name, users.picture, news.upvote, news.downvote, 
		news.comment, news.view, TO_CHAR(news.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM news
		INNER JOIN categories ON news.category_id = categories.id
		INNER JOIN users ON news.author_id = users.id
		LEFT JOIN news_images ON news.id = news_images.news_id
		WHERE news.id = $1;
	`, id)

	count := 0
	for rows.Next() {
		var newImage string
		var errScan error

		if count < 1 {
			errScan = rows.Scan(
				&news.Id, &news.Title, &news.Content, &news.SlugUrl, &news.CoverImage,
				&newImage, &news.Nsfw, &category.Id, &category.Name, &author.Id,
				&author.Name, &author.Picture, &counter.Upvote, &counter.Downvote, &counter.Comment,
				&counter.View, &news.CreatedAt,
			)
		} else {
			errScan = rows.Scan(nil, nil, nil, nil, nil, &newImage, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		}

		if errScan != nil {
			return news, errScan
		}

		if newImage != "" {
			newsImage = append(newsImage, newImage)
		}

		count++
	}

	news.AdditionalImages = newsImage
	news.Category = category
	news.Author = author
	news.Counter = counter

	return news, err
}

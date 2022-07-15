package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type NewsRepository struct {
	DB *pgx.Conn
}

func NewNewsRepository(DB *pgx.Conn) NewsRepository {
	return NewsRepository{DB}
}

func (r NewsRepository) SelectAll(category string, scope string) (pgx.Rows, error) {
	var rows pgx.Rows
	var err error

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

	return rows, err
}

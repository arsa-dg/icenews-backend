package repository

import (
	"context"
	"database/sql"
	"icenews/backend/model"

	"github.com/google/uuid"
)

type NewsRepositoryInterface interface {
	SelectAll(category string, scope string) ([]model.NewsListRaw, error)
	SelectById(id string) ([]model.NewsDetailRaw, error)
	SelectAllCategory() ([]model.NewsCategory, error)
	InsertComment(description, newsId string, authorId uuid.UUID) (int, error)
	SelectCommentByNewsId(newsId string) ([]model.Comment, error)
}

type NewsRepository struct {
	DB *sql.DB
}

func NewNewsRepository(DB *sql.DB) NewsRepository {
	return NewsRepository{DB}
}

func (r NewsRepository) SelectAll(category string, scope string) ([]model.NewsListRaw, error) {
	var rows *sql.Rows
	var err error
	var newsListRaw []model.NewsListRaw

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
		rows, err = r.DB.QueryContext(context.Background(), query+` WHERE
			news.category_id = $1 AND
			news.scope = $2;
		`, category, scope)
	} else if category != "" {
		rows, err = r.DB.QueryContext(context.Background(), query+` WHERE
			news.category_id = $1;
		`, category)
	} else if scope != "" {
		rows, err = r.DB.QueryContext(context.Background(), query+` WHERE
			news.scope = $1;
		`, scope)
	} else {
		rows, err = r.DB.QueryContext(context.Background(), query)
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		news := model.NewsListRaw{}

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

func (r NewsRepository) SelectById(id string) ([]model.NewsDetailRaw, error) {
	newsDetailRaw := []model.NewsDetailRaw{}

	rows, err := r.DB.QueryContext(context.Background(), `SELECT
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

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		news := model.NewsDetailRaw{}

		errScan := rows.Scan(
			&news.Id, &news.Title, &news.Content, &news.SlugUrl, &news.CoverImage,
			&news.AdditionalImage, &news.Nsfw, &news.CategoryId,
			&news.CategoryName, &news.AuthorId, &news.AuthorName,
			&news.AuthorPicture, &news.Upvote, &news.Downvote,
			&news.Comment, &news.View, &news.CreatedAt,
		)

		if errScan != nil {
			return nil, errScan
		}

		newsDetailRaw = append(newsDetailRaw, news)
	}

	return newsDetailRaw, err
}

func (r NewsRepository) SelectAllCategory() ([]model.NewsCategory, error) {
	rows, err := r.DB.QueryContext(context.Background(), "SELECT * FROM categories;")

	if err != nil {
		return nil, err
	}

	categoryList := []model.NewsCategory{}

	for rows.Next() {
		category := model.NewsCategory{}

		errScan := rows.Scan(&category.Id, &category.Name)

		if errScan != nil {
			return nil, errScan
		}

		categoryList = append(categoryList, category)
	}

	return categoryList, nil
}

func (r NewsRepository) InsertComment(description, newsId string, authorId uuid.UUID) (int, error) {
	var commentId int

	err := r.DB.QueryRowContext(context.Background(), "INSERT INTO comments(description, author_id, news_id) values($1, $2, $3) RETURNING id;",
		description, authorId, newsId).Scan(&commentId)

	if err != nil {
		return -1, err
	}

	return commentId, err
}

func (r NewsRepository) SelectCommentByNewsId(newsId string) ([]model.Comment, error) {
	rows, err := r.DB.QueryContext(context.Background(), `SELECT
		comments.id, comments.description, users.id, users.name, users.picture, 
		TO_CHAR(comments.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z')
		FROM comments 
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.news_id = $1;
	`, newsId)

	if err != nil {
		return nil, err
	}

	commentList := []model.Comment{}

	for rows.Next() {
		comment := model.Comment{}

		errScan := rows.Scan(
			&comment.Id, &comment.Description, &comment.Commentator.Id,
			&comment.Commentator.Name, &comment.Commentator.Picture,
			&comment.CreatedAt,
		)

		if errScan != nil {
			return nil, errScan
		}

		commentList = append(commentList, comment)
	}

	return commentList, nil
}

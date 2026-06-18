package repository

import (
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
)

type PostRepository struct {
	db *PostgresDB
}

func NewPostRepository(db *PostgresDB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (id, author_id, content, media_url, created_at) VALUES ($1,$2, $3,$4, $5)`
	_, err := r.db.Pool.Exec(ctx, query, post.ID, post.AuthorID, post.Content, post.MediaURL, post.CreatedAt)
	return err
}

func (r *PostRepository) GetFeed(ctx context.Context, limit, offset int, currentUserID uuid.UUID) ([]*models.Post, error) {
	// Убедись, что здесь четкие пробелы вокруг 1и1 и1и2
	query := `SELECT p.id, p.author_id, u.username, p.content, p.media_url, p.created_at, 
       (SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) as likes_count,
            EXISTS (SELECT 1 FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) as is_liked
              FROM posts p
              JOIN users u ON p.author_id = u.id
              ORDER BY p.created_at DESC 
              LIMIT $2 OFFSET $3`

	rows, err := r.db.Pool.Query(ctx, query, currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.AuthorName, &p.Content, &p.MediaURL, &p.CreatedAt, &p.LikesCount, &p.IsLiked); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetFeedByAuthors(ctx context.Context, authorIDs []uuid.UUID, limit, offset int) ([]*models.Post, error) {
	if len(authorIDs) == 0 {
		return []*models.Post{}, nil // Если ни на кого не подписан, лента пуста
	}

	// Используем оператор ANY($1), чтобы передать массив ID в Postgres
	query := `SELECT id, author_id, content, media_url, created_at 
              FROM posts 
              WHERE author_id = ANY($1) 
              ORDER BY created_at DESC 
              LIMIT $2 OFFSET $3`

	rows, err := r.db.Pool.Query(ctx, query, authorIDs, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		p := &models.Post{}
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

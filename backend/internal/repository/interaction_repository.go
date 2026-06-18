package repository

import (
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
)

type InteractionRepository struct {
	db *PostgresDB
}

func NewInteractionRepository(db *PostgresDB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

func (r *InteractionRepository) LikePost(ctx context.Context, userID, postID uuid.UUID) error {
	query := `INSERT INTO likes (post_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`
	_, err := r.db.Pool.Exec(ctx, query, postID, userID)
	return err
}

func (r *InteractionRepository) UnlikePost(ctx context.Context, userID, postID uuid.UUID) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND post_id = $2`
	ct, err := r.db.Pool.Exec(ctx, query, postID, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func (r *InteractionRepository) AddComment(ctx context.Context, authorID, postID uuid.UUID, content string) error {
	query := `INSERT INTO comments (id, post_id, author_id, content, created_at) VALUES ($1, $2, $3,$4, NOW())`
	_, err := r.db.Pool.Exec(ctx, query, uuid.New(), postID, authorID, content)
	return err
}

func (r *InteractionRepository) GetCommentsByPost(ctx context.Context, postID uuid.UUID) ([]*models.Comment, error) {
	query := `
        SELECT c.id, c.post_id, c.author_id, u.username, c.content, c.created_at 
        FROM comments c
        JOIN users u ON c.author_id = u.id
        WHERE c.post_id = $1
        ORDER BY c.created_at ASC` // Старые комментарии сверху, новые снизу

	rows, err := r.db.Pool.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.AuthorName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

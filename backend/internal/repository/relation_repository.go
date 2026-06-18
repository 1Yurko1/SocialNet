package repository

import (
	"context"

	"github.com/google/uuid"
)

type RelationRepository struct {
	db *PostgresDB
}

func NewRelationRepository(db *PostgresDB) *RelationRepository {
	return &RelationRepository{db: db}
}

func (r *RelationRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `INSERT INTO relations (follower_id, following_id) VALUES ($1,$2) 
              ON CONFLICT DO NOTHING` // Чтобы не было ошибки при повторной подписке
	_, err := r.db.Pool.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *RelationRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	query := `DELETE FROM relations WHERE follower_id = $1 AND following_id = $2`

	_, err := r.db.Pool.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *RelationRepository) GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT following_id FROM relations WHERE follower_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"strings"
)

type UserRepository struct {
	db *PostgresDB
}

func NewUserRepository(db *PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, email, password, created_at, updated_at) 
              VALUES ($1,$2, $3,$4, $5,$6)`

	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password FROM users WHERE username = $1`

	err := r.db.Pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password FROM users WHERE email = $1`

	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) SearchUsers(ctx context.Context, query string) ([]*models.User, error) {
	sqlQuery := `SELECT id, username, COALESCE(avatar_url, '') 
	             FROM users 
	             WHERE username ILIKE $1 
	             LIMIT 10`

	// Экранируем спецсимволы LIKE и формируем паттерн безопасно
	pattern := "%" + strings.ReplaceAll(strings.ReplaceAll(query, "%", `\%`), "_", `\_`) + "%"

	rows, err := r.db.Pool.Query(ctx, sqlQuery, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// ✅ Обязательная проверка ошибки после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate search results: %w", err)
	}

	return users, nil
}

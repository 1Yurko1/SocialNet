package repository

import (
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ChatRepository struct {
	db *PostgresDB
}

func NewChatRepository(db *PostgresDB) *ChatRepository {
	return &ChatRepository{db: db}
}

// GetOrCreatePrivateChat создает приватный чат между двумя людьми или возвращает существующий
func (r *ChatRepository) GetOrCreatePrivateChat(ctx context.Context, user1, user2 uuid.UUID) (uuid.UUID, error) {
	var chatID uuid.UUID

	// Ищем существующий приватный чат между двумя пользователями
	query := `SELECT c.id FROM chats c
          JOIN chat_members m1 ON c.id = m1.chat_id
          JOIN chat_members m2 ON c.id = m2.chat_id
          WHERE c.type = 'private'
            AND m1.user_id = $1
            AND m2.user_id = $2
            AND c.id NOT IN (
                SELECT chat_id FROM chat_members 
                GROUP BY chat_id 
                HAVING COUNT(*) > 2
            )`

	err := r.db.Pool.QueryRow(ctx, query, user1, user2).Scan(&chatID)
	if err == nil {
		return chatID, nil
	}
	if err != pgx.ErrNoRows {
		return uuid.Nil, err
	}

	// Чата нет — создаём в транзакции
	chatID = uuid.New()
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Создаём чат
	_, err = tx.Exec(ctx, "INSERT INTO chats (id, type) VALUES ($1, 'private')", chatID)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. Добавляем обоих участников
	_, err = tx.Exec(ctx,
		"INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)",
		chatID, user1, user2,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return chatID, nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *models.Message) error {
	query := `INSERT INTO messages (id, chat_id, sender_id, content, created_at)
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Pool.Exec(ctx, query,
		msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.CreatedAt,
	)
	return err
}

func (r *ChatRepository) GetMessagesByChat(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	query := `SELECT id, chat_id, sender_id, content, created_at 
              FROM messages 
              WHERE chat_id = $1 
              ORDER BY created_at ASC 
              LIMIT $2 OFFSET $3`

	rows, err := r.db.Pool.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*models.Message, 0)

	for rows.Next() {
		m := &models.Message{}
		if err := rows.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	// Проверяем, не случилось ли ошибки при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return messages, nil
}

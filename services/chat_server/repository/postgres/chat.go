package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Genvekt/cli-chat/services/chat-server/repository"

	"github.com/Genvekt/cli-chat/services/chat-server/model"
)

const (
	chatTable = "chat"
	idColumn  = "id"
)

var _ repository.ChatRepository = (*ChatPostgresRepository)(nil)

// ChatPostgresRepository implements repository.ChatRepository for postgres data source
type ChatPostgresRepository struct {
	db *pgxpool.Pool
}

// NewChatPostgresRepository retrieves new ChatPostgresRepository instance
func NewChatPostgresRepository(db *pgxpool.Pool) *ChatPostgresRepository {
	return &ChatPostgresRepository{
		db: db,
	}
}

// Create inserts chat into db
func (r *ChatPostgresRepository) Create(ctx context.Context, chat *model.Chat) error {
	builderInsert := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&chat.ID)
	if err != nil {
		return fmt.Errorf("failed to insert chat: %v", err)
	}

	return nil
}

// Delete deletes chat by id
func (r *ChatPostgresRepository) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete chat with id %d: %v", id, err)
	}

	return nil
}

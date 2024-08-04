package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat/converter"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

const (
	chatTable      = "chat"
	chatIDColumn   = "id"
	chatNameColumn = "name"
)

var _ repository.ChatRepository = (*chatPostgresRepository)(nil)

// ChatPostgresRepository implements repository.ChatRepository for postgres data source
type chatPostgresRepository struct {
	db db.Client
}

// NewChatPostgresRepository retrieves new chatPostgresRepository instance
func NewChatPostgresRepository(db db.Client) *chatPostgresRepository {
	return &chatPostgresRepository{
		db: db,
	}
}

// Create inserts chat into db
func (r *chatPostgresRepository) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	chatDB := converter.ToRepoFromChat(chat)

	builderChatInsert := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatNameColumn).
		Values(chatDB.Name).
		Suffix("RETURNING id")

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.CreateChat",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatDB.ID)
	if err != nil {
		return 0, err
	}

	return chatDB.ID, nil
}

// Delete deletes chat by id
func (r *chatPostgresRepository) Delete(ctx context.Context, id int64) error {
	builderChatDelete := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: id})

	query, args, err := builderChatDelete.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.DeleteChat",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

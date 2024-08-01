package message

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/db"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository/message/converter"
)

const (
	messageTable           = "message"
	messageSenderIDColumn  = "sender_id"
	messageChatIDColumn    = "chat_id"
	messageContentColumn   = "content"
	messageTimestampColumn = "sent_at"
)

var _ repository.MessageRepository = (*messagePostgresRepository)(nil)

type messagePostgresRepository struct {
	db db.Client
}

// NewMessagePostgresRepository initialises repository of messages
func NewMessagePostgresRepository(db db.Client) *messagePostgresRepository {
	return &messagePostgresRepository{
		db: db,
	}
}

// Create creates messages in db
func (r *messagePostgresRepository) Create(ctx context.Context, message *model.Message) error {
	messageDB := converter.ToRepoFromMessage(message)

	builderChatInsert := sq.Insert(messageTable).
		PlaceholderFormat(sq.Dollar).
		Columns(messageSenderIDColumn, messageChatIDColumn, messageContentColumn, messageTimestampColumn).
		Values(messageDB.SenderID, messageDB.ChatID, messageDB.Content, messageDB.Timestamp)

	query, args, err := builderChatInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "message_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	return nil
}

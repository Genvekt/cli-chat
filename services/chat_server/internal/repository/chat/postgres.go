package chat

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/db"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat/converter"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat/model"
)

const (
	chatTable      = "chat"
	chatIDColumn   = "id"
	chatNameColumn = "name"

	chatMemberTable          = "chat_member"
	chatMemberUserIDColumn   = "user_id"
	chatMemberChatIDColumn   = "chat_id"
	chatMemberJoinedAtColumn = "joined_at"
)

var _ repository.ChatRepository = (*chatPostgresRepository)(nil)

// ChatPostgresRepository implements repository.ChatRepository for postgres data source
type chatPostgresRepository struct {
	db db.Client
}

// NewChatPostgresRepository retrieves new ChatPostgresRepository instance
func NewChatPostgresRepository(db db.Client) *chatPostgresRepository {
	return &chatPostgresRepository{
		db: db,
	}
}

// Create inserts chat into db
func (r *chatPostgresRepository) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	chatDB, chatMembersDB := converter.ToRepoFromChat(chat)
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
		return 0, fmt.Errorf("failed to insert chat: %v", err)
	}

	builderMembersInsert := sq.Insert(chatMemberTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatMemberUserIDColumn, chatMemberChatIDColumn, chatMemberJoinedAtColumn)

	for _, member := range chatMembersDB {
		builderMembersInsert = builderMembersInsert.Values(member.UserID, chatDB.ID, member.JoinedAt)
	}

	query, args, err = builderMembersInsert.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	q = db.Query{
		Name:     "chat_repository.CreateChatMember",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert chat members: %v", err)
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
		return fmt.Errorf("failed to delete chat with id %d: %v", id, err)
	}

	return nil
}

// GetChatMember retrieves chat member
func (r *chatPostgresRepository) GetChatMember(
	ctx context.Context,
	chatID int64,
	userID int64,
) (*model.ChatMember, error) {
	builderSelect := sq.Select(chatMemberChatIDColumn, chatMemberUserIDColumn, chatMemberJoinedAtColumn).
		From(chatMemberTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{
			sq.Eq{chatMemberChatIDColumn: chatID},
			sq.Eq{chatMemberUserIDColumn: userID},
		}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.GetChatMember",
		QueryRaw: query,
	}

	member := &repoModel.Member{}

	err = r.db.DB().ScanOneContext(ctx, member, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrChatMemberNotFound
		}
		return nil, fmt.Errorf("failed to get chat member: %v", err)
	}

	return converter.ToChatMemberFromRepo(member), nil
}

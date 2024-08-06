package chat_member

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/db"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat_member/converter"
	repoModel "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat_member/model"
)

const (
	chatMemberTable          = "chat_member"
	chatMemberUserIDColumn   = "user_id"
	chatMemberChatIDColumn   = "chat_id"
	chatMemberJoinedAtColumn = "joined_at"
)

var _ repository.ChatMemberRepository = (*chatMemberPostgresRepository)(nil)

type chatMemberPostgresRepository struct {
	db db.Client
}

// NewChatMemberPostgresRepository retrieves new chatMemberPostgresRepository instance
func NewChatMemberPostgresRepository(db db.Client) *chatMemberPostgresRepository {
	return &chatMemberPostgresRepository{
		db: db,
	}
}

// CreateBatch creates many chat members at once
func (r *chatMemberPostgresRepository) CreateBatch(ctx context.Context, chatID int64, userIDs []int64) error {
	builderMembersInsert := sq.Insert(chatMemberTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatMemberChatIDColumn, chatMemberUserIDColumn)

	for _, userID := range userIDs {
		builderMembersInsert = builderMembersInsert.Values(chatID, userID)
	}

	query, args, err := builderMembersInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_member_repository.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves member of a chat by chat id and user id
func (r *chatMemberPostgresRepository) Get(
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
		Name:     "chat_member_repository.Get",
		QueryRaw: query,
	}

	member := &repoModel.Member{}

	err = r.db.DB().ScanOneContext(ctx, member, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrChatMemberNotFound
		}
		return nil, err
	}

	return converter.ToChatMemberFromRepo(member), nil
}

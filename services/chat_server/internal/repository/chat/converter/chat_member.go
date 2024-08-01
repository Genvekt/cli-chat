package converter

import (
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat/model"
)

// ToChatMemberFromRepo repository model chat member to service model chat member
func ToChatMemberFromRepo(chatMember *repoModel.Member) *model.ChatMember {
	return &model.ChatMember{
		ID:       chatMember.UserID,
		JoinedAt: chatMember.JoinedAt,
	}
}

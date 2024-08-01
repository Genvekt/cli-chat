package converter

import (
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat/model"
)

// ToChatFromRepo converts repository model chat to service model chat
func ToChatFromRepo(chat *repoModel.Chat, members []*repoModel.Member) *model.Chat {
	chatMembers := make([]*model.ChatMember, len(members))
	for _, member := range members {
		chatMembers = append(chatMembers, ToChatMemberFromRepo(member))
	}
	return &model.Chat{
		ID:        chat.ID,
		Name:      chat.Name,
		Members:   chatMembers,
		CreatedAt: chat.CreatedAt,
	}
}

// ToRepoFromChat converts service model chat to repository model chat
func ToRepoFromChat(chat *model.Chat) (*repoModel.Chat, []*repoModel.Member) {
	chatRepo := &repoModel.Chat{
		ID:        chat.ID,
		Name:      chat.Name,
		CreatedAt: chat.CreatedAt,
	}

	chatMembers := make([]*repoModel.Member, 0, len(chat.Members))
	for _, member := range chat.Members {
		chatMembers = append(chatMembers, &repoModel.Member{
			UserID:   member.ID,
			ChatID:   chat.ID,
			JoinedAt: member.JoinedAt,
		})
	}

	return chatRepo, chatMembers
}

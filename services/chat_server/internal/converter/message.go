package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// ToMessageFromProto converts message repository model to message service model
func ToMessageFromProto(message *chatApi.Message) *model.Message {
	return &model.Message{
		SenderID:  message.SenderId,
		ChatID:    message.ChatId,
		Content:   message.Text,
		Timestamp: message.Timestamp.AsTime(),
	}
}

// ToProtoFromMessage converts message service model to message repository model
func ToProtoFromMessage(message *model.Message) *chatApi.Message {
	return &chatApi.Message{
		SenderId:  message.SenderID,
		ChatId:    message.ChatID,
		Text:      message.Content,
		Timestamp: timestamppb.New(message.Timestamp),
	}
}

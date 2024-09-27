package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/model"
)

func ToProtoFromMessage(message *model.Message) *chatApi.Message {
	return &chatApi.Message{
		SenderId:  message.SenderID,
		ChatId:    message.ChatID,
		Text:      message.Text,
		Timestamp: timestamppb.New(message.Timestamp),
	}
}

func FromProtoToMessage(protoMessage *chatApi.Message) *model.Message {
	return &model.Message{
		SenderID:  protoMessage.SenderId,
		ChatID:    protoMessage.ChatId,
		Text:      protoMessage.Text,
		Timestamp: protoMessage.Timestamp.AsTime(),
	}
}

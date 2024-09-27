package chat

import (
	"go.uber.org/zap"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/converter"
)

// ConnectChat creates stream connection to chat
func (s *Service) ConnectChat(req *chatApi.ConnectChatRequest, stream chatApi.ChatV1_ConnectChatServer) error {
	messageChan, err := s.chatService.Connect(stream.Context(), req.GetChatId(), req.GetUsername())
	if err != nil {
		return err
	}

	logger.Debug("user connected to chat",
		zap.Int64("chat_id", req.GetChatId()),
		zap.String("username", req.GetUsername()),
	)

	defer logger.Debug("user disconnected from chat",
		zap.Int64("chat_id", req.GetChatId()),
		zap.String("username", req.GetUsername()),
	)

	for {
		select {

		case msg, ok := <-messageChan:
			if !ok {
				return nil
			}
			if err = stream.Send(converter.ToProtoFromMessage(msg)); err != nil {
				return err
			}

		case <-stream.Context().Done():
			return nil
		}
	}

}

package chat

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// Connect initialises the connection of user to chat
func (s *chatService) Connect(ctx context.Context, id int64, username string) (chan *model.Message, error) {
	s.mxChat.Lock()
	chat, isActive := s.chatConnections[id]
	if !isActive {
		chat = &model.ChatConnection{
			ChatID:      id,
			Buffer:      make(chan *model.Message, 100),
			Connections: map[string]chan *model.Message{},
		}

		// Start chat broadcasting logic
		go s.processChatBuffer(chat)

		s.chatConnections[id] = chat

		logger.Debug("chat activated", zap.Int64("chat_id", id))
	}
	s.mxChat.Unlock()

	chat.Mx.Lock()
	defer chat.Mx.Unlock()

	if _, connected := chat.Connections[username]; connected {
		return nil, fmt.Errorf("already connected")
	}

	userConnection := make(chan *model.Message, 100)

	chat.Connections[username] = userConnection

	go s.handleDisconnect(ctx, id, username)

	return userConnection, nil
}

func (s *chatService) handleDisconnect(ctx context.Context, chatID int64, username string) {
	<-ctx.Done()
	s.mxChat.RLock()
	chat, isActive := s.chatConnections[chatID]
	s.mxChat.RUnlock()

	if !isActive {
		return
	}

	chat.Mx.Lock()
	delete(chat.Connections, username)
	chat.Mx.Unlock()

	if len(chat.Connections) == 0 {
		// deactivate chat if there are no connected users
		chat.Mx.Lock()
		close(chat.Buffer)
		chat.Mx.Unlock()

		s.mxChat.Lock()
		delete(s.chatConnections, chatID)
		s.mxChat.Unlock()

		logger.Debug("chat deactivated", zap.Int64("chat_id", chatID))
	}
}

func (s *chatService) processChatBuffer(chatConnection *model.ChatConnection) {
	for {
		msg, ok := <-chatConnection.Buffer
		if !ok {
			logger.Debug("chat buffer closed", zap.Int64("chat_id", chatConnection.ChatID))
			return
		}

		chatConnection.Mx.RLock()
		for _, connection := range chatConnection.Connections {
			connection <- msg
		}
		chatConnection.Mx.RUnlock()
	}
}

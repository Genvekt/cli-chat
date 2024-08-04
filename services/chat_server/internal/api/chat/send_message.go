package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/converter"
)

// SendMessage handles SendMessageRequest
func (s *Service) SendMessage(ctx context.Context, req *chatApi.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.chatService.SendMessage(ctx, converter.ToMessageFromProto(req.Message))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

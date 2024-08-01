package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
)

// Delete handles DeleteRequest
func (s *Service) Delete(ctx context.Context, req *chatApi.DeleteRequest) (*emptypb.Empty, error) {
	err := s.chatService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

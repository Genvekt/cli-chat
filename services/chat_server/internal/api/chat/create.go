package chat

import (
	"context"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
)

// Create handles CreateRequest
func (s *Service) Create(ctx context.Context, req *chatApi.CreateRequest) (*chatApi.CreateResponse, error) {
	newChatID, err := s.chatService.Create(ctx, req.Name, req.Usernames)
	if err != nil {
		return nil, err
	}

	return &chatApi.CreateResponse{
		Id: newChatID,
	}, nil
}

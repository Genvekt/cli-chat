package chat

import (
	"context"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
)

func validateCreateRequest(req *chatApi.CreateRequest) error {
	if len(req.Usernames) == 0 {
		return ErrEmptyChat
	}
	return nil
}

// Create handles CreateRequest
func (s *Service) Create(ctx context.Context, req *chatApi.CreateRequest) (*chatApi.CreateResponse, error) {
	err := validateCreateRequest(req)
	if err != nil {
		return nil, err
	}

	newChatID, err := s.chatService.Create(ctx, req.Name, req.Usernames)
	if err != nil {
		return nil, err
	}

	return &chatApi.CreateResponse{
		Id: newChatID,
	}, nil
}

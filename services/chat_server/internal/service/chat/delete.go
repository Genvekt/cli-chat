package chat

import (
	"context"
	"fmt"
)

// Delete deletes chat
func (s *chatService) Delete(ctx context.Context, id int64) error {
	err := s.chatRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot delete chat: %w", err)
	}

	return nil
}

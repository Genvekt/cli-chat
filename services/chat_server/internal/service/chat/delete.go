package chat

import "context"

// Delete deletes chat
func (s *chatService) Delete(ctx context.Context, id int64) error {
	err := s.chatRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

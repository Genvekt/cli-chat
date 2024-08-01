package user

import "context"

// Delete deletes user by id
func (s *userService) Delete(ctx context.Context, id int64) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

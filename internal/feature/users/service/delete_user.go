package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(ctx context.Context, id int) error {
	if err := s.userRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("error to delete user from repository:%w", err)
	}
	return nil
}

package users_service

import (
	"context"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

func (s *UsersService) PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.Users, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return domain.Users{}, fmt.Errorf("get user: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.Users{}, fmt.Errorf("apply user patch: %w", err)
	}

	pathedUser, err := s.userRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.Users{}, fmt.Errorf("patched user: %w", err)
	}

	return pathedUser, nil

}

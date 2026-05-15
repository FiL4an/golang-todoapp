package users_service

import (
	"context"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_errors "github.com/FiL4an/golang-todoapp/internal/core/errors"
)

func (s *UsersService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.Users, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative: %w", core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative: %w ", core_errors.ErrInvalidArgument)
	}
	usersDomain, err := s.userRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get users from  repository: %w", err)
	}

	return usersDomain, nil
}

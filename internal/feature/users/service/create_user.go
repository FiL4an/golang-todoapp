package users_service

import (
	"context"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

func (s *UsersService) CreateUser(ctx context.Context, user domain.Users) (domain.Users, error) {
	//1. uss`t.validate
	//2.

	if err := user.Validate(); err != nil {
		return domain.Users{}, fmt.Errorf("validate user domain: %w ", err)
	}
	user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.Users{}, fmt.Errorf("create user:%w ", err)
	}
	return user, nil
}

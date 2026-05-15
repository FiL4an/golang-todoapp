package users_service

import (
	"context"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
)

type UsersService struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.Users) (domain.Users, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.Users, error)
	GetUser(ctx context.Context, id int) (domain.Users, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, user domain.Users) (domain.Users, error)
}

func NewUsersService(userRepository UserRepository) *UsersService {
	return &UsersService{
		userRepository: userRepository,
	}
}

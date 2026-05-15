package users_postgres_repository

import "github.com/FiL4an/golang-todoapp/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func UserDomainsFromModels(users []UserModel) []domain.Users {
	userDomains := make([]domain.Users, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUsers(user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber)

	}
	return userDomains
}

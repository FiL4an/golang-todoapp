package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/FiL4an/golang-todoapp/internal/core/domain"
	core_errors "github.com/FiL4an/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/FiL4an/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UserRepository) PatchUser(ctx context.Context, id int, user domain.Users) (domain.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET full_name= $1,
		phone_number=$2, 
		version = version+1
	
	WHERE id=$3 AND version=$4
	RETURNING 
	id, 
	version, 
	full_name, 
	phone_number;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		user.ID,
		user.Version,
	)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Users{}, fmt.Errorf(
				"user with  id='%d' concurrently accessed: %w",
				id, core_errors.ErrConflict)
		} else {
			return domain.Users{}, fmt.Errorf("scan error : %w ", err)
		}
	}

	userDomain := domain.NewUsers(userModel.ID, userModel.Version, user.FullName, user.PhoneNumber)

	return userDomain, nil
}

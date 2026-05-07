package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/FiL4an/golang-todoapp/internal/core/errors"
)

type Users struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUsers(id int, version int, fullName string, phoneNumber *string) Users {
	return Users{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}
func NewUserUnItialazed(fullName string, phoneNumber *string) Users {
	return NewUsers(UninitializedID, UninitializedVersion, fullName, phoneNumber)

}

func (u *Users) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid `FullName` len:%d:%w", fullNameLength, core_errors.ErrInvalidArgument)
	}
	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("invalide `PhoneNumber` len:%d:%w", phoneNumberLen, core_errors.ErrInvalidArgument)
		}
		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("Invalid `PhoneNumber` format:%w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (u *UserPatch) Validate() error {
	if u.FullName.Set && u.FullName.Value == nil {
		return fmt.Errorf(
			"FullName can't be patched to NULL: %w ",
			core_errors.ErrInvalidArgument)
	}
	return nil
}

func (u *Users) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch:%w", err)
	}

	tmp := *u
	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w ", err)
	}
	*u = tmp
	return nil
}

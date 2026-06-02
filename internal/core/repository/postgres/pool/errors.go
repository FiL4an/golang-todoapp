package core_postgres_pool

import "errors"

var (
	ErrNoRows              = errors.New("no rows")
	ErrViolatesForgivenKey = errors.New("violetes foreign key")
	ErrUnknown             = errors.New("unknown")
)

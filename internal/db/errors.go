package db

import "errors"

var (
	Err                 = errors.New("placeholder")
	ErrUserDoesNotExist = errors.New("provided user does not exist")
	ErrInactiveUser     = errors.New("user is inactive")
)

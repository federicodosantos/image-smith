package error

import "errors"

var (
	ErrEmailNotFound     = errors.New("email not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailExist        = errors.New("email already exist")
	ErrNotVerified       = errors.New("account has not been verified")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrDatabase          = errors.New("database error")
	ErrRowsAffected      = errors.New("error due to there is no or more than 1 affected column")
)

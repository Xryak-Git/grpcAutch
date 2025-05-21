package services

import "errors"

var (
	ErrInvalidCredentials = errors.New("login or password is incorrect")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("user already exists")
)

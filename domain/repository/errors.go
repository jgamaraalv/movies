package repository

import "errors"

var (
	ErrMovieNotFound            = errors.New("movie not found")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrRegistrationValidation   = errors.New("registration failed")
)

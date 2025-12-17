package repository

import "errors"

// Repository errors (persistence)
var (
	ErrMovieNotFound = errors.New("movie not found")
	ErrUserNotFound  = errors.New("user not found")
)

// Domain errors (business rules)
var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrRegistrationValidation   = errors.New("registration failed")
	ErrNameRequired             = errors.New("name is required")
)

// Collection errors
var (
	ErrMovieAlreadyInFavorites = errors.New("movie already in favorites")
	ErrMovieAlreadyInWatchlist = errors.New("movie already in watchlist")
	ErrMovieNotInFavorites     = errors.New("movie not in favorites")
	ErrMovieNotInWatchlist     = errors.New("movie not in watchlist")
	ErrInvalidCollectionType   = errors.New("invalid collection type")
)

package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
)

const (
	CollectionFavorite  = "favorite"
	CollectionWatchlist = "watchlist"
)

type SaveToCollectionInput struct {
	Email      string
	MovieID    int
	Collection string
}

type SaveToCollectionOutput struct {
	Success bool
	Message string
}

type SaveToCollectionUseCase struct {
	accountStorage providers.AccountStorage
	logger         *logger.Logger
}

func NewSaveToCollectionUseCase(storage providers.AccountStorage, log *logger.Logger) *SaveToCollectionUseCase {
	return &SaveToCollectionUseCase{
		accountStorage: storage,
		logger:         log,
	}
}

func (uc *SaveToCollectionUseCase) Execute(input SaveToCollectionInput) (*SaveToCollectionOutput, error) {
	if err := uc.validateInput(input); err != nil {
		return nil, err
	}

	user := models.User{Email: input.Email}
	success, err := uc.accountStorage.SaveCollection(user, input.MovieID, input.Collection)
	if err != nil {
		uc.logger.Error("Failed to save movie to collection", err)
		return nil, err
	}

	message := "Movie added to " + input.Collection + " successfully"
	uc.logger.Info(message + " for user: " + input.Email)

	return &SaveToCollectionOutput{
		Success: success,
		Message: message,
	}, nil
}

func (uc *SaveToCollectionUseCase) validateInput(input SaveToCollectionInput) error {
	if input.Email == "" {
		return errors.New("email is required")
	}
	if input.MovieID <= 0 {
		return errors.New("invalid movie ID")
	}
	if input.Collection != CollectionFavorite && input.Collection != CollectionWatchlist {
		return errors.New("collection must be 'favorite' or 'watchlist'")
	}
	return nil
}

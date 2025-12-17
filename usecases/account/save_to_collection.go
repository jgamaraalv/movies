package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/domain/valueobject"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
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
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewSaveToCollectionUseCase(repo repository.UserRepository, log *logger.Logger) *SaveToCollectionUseCase {
	return &SaveToCollectionUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *SaveToCollectionUseCase) Execute(input SaveToCollectionInput) (*SaveToCollectionOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if input.MovieID <= 0 {
		return nil, errors.New("invalid movie ID")
	}

	if input.Collection != CollectionFavorite && input.Collection != CollectionWatchlist {
		return nil, errors.New("collection must be 'favorite' or 'watchlist'")
	}

	user := models.User{Email: email.String()}
	success, err := uc.userRepo.SaveCollection(user, input.MovieID, input.Collection)
	if err != nil {
		uc.logger.Error("Failed to save movie to collection", err)
		return nil, err
	}

	message := "Movie added to " + input.Collection + " successfully"
	uc.logger.Info(message + " for user: " + email.String())

	return &SaveToCollectionOutput{
		Success: success,
		Message: message,
	}, nil
}

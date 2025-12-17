package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/domain/entity"
	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/domain/valueobject"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type SaveToCollectionInput struct {
	Email      string
	MovieID    int
	Collection string
}

type SaveToCollectionOutput struct {
	Success             bool
	Message             string
	AlreadyInCollection bool
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

	if input.Collection != entity.CollectionFavorites && input.Collection != entity.CollectionWatchlist {
		return nil, repository.ErrInvalidCollectionType
	}

	userModel, err := uc.userRepo.GetAccountDetails(email.String())
	if err != nil {
		uc.logger.Error("Failed to get user details", err)
		return nil, err
	}

	user, err := entity.UserFromModel(userModel)
	if err != nil {
		uc.logger.Error("Failed to convert user model to entity", err)
		return nil, err
	}

	if user.IsInCollection(input.MovieID, input.Collection) {
		collectionName := "favorites"
		if input.Collection == entity.CollectionWatchlist {
			collectionName = "watchlist"
		}
		return &SaveToCollectionOutput{
			Success:             true,
			Message:             "Movie is already in " + collectionName,
			AlreadyInCollection: true,
		}, nil
	}

	userModelForSave := models.User{Email: email.String()}
	success, err := uc.userRepo.SaveCollection(userModelForSave, input.MovieID, input.Collection)
	if err != nil {
		uc.logger.Error("Failed to save movie to collection", err)
		return nil, err
	}

	message := "Movie added to " + input.Collection + " successfully"
	uc.logger.Info(message + " for user: " + email.String())

	return &SaveToCollectionOutput{
		Success:             success,
		Message:             message,
		AlreadyInCollection: false,
	}, nil
}

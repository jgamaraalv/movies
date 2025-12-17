package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
)

type GetFavoritesInput struct {
	Email string
}

type GetFavoritesOutput struct {
	Favorites []models.Movie
}

type GetFavoritesUseCase struct {
	accountStorage providers.AccountStorage
	logger         *logger.Logger
}

func NewGetFavoritesUseCase(storage providers.AccountStorage, log *logger.Logger) *GetFavoritesUseCase {
	return &GetFavoritesUseCase{
		accountStorage: storage,
		logger:         log,
	}
}

func (uc *GetFavoritesUseCase) Execute(input GetFavoritesInput) (*GetFavoritesOutput, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}

	user, err := uc.accountStorage.GetAccountDetails(input.Email)
	if err != nil {
		uc.logger.Error("Failed to get user favorites", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved favorites for: " + input.Email)

	return &GetFavoritesOutput{
		Favorites: user.Favorites,
	}, nil
}

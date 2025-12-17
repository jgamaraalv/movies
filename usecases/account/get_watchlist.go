package account

import (
	"errors"

	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
)

type GetWatchlistInput struct {
	Email string
}

type GetWatchlistOutput struct {
	Watchlist []models.Movie
}

type GetWatchlistUseCase struct {
	accountStorage providers.AccountStorage
	logger         *logger.Logger
}

func NewGetWatchlistUseCase(storage providers.AccountStorage, log *logger.Logger) *GetWatchlistUseCase {
	return &GetWatchlistUseCase{
		accountStorage: storage,
		logger:         log,
	}
}

func (uc *GetWatchlistUseCase) Execute(input GetWatchlistInput) (*GetWatchlistOutput, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}

	user, err := uc.accountStorage.GetAccountDetails(input.Email)
	if err != nil {
		uc.logger.Error("Failed to get user watchlist", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved watchlist for: " + input.Email)

	return &GetWatchlistOutput{
		Watchlist: user.Watchlist,
	}, nil
}

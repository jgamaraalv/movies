package account

import (
	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/domain/valueobject"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type GetWatchlistInput struct {
	Email string
}

type GetWatchlistOutput struct {
	Watchlist []models.Movie
}

type GetWatchlistUseCase struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewGetWatchlistUseCase(repo repository.UserRepository, log *logger.Logger) *GetWatchlistUseCase {
	return &GetWatchlistUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *GetWatchlistUseCase) Execute(input GetWatchlistInput) (*GetWatchlistOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetAccountDetails(email.String())
	if err != nil {
		uc.logger.Error("Failed to get user watchlist", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved watchlist for: " + email.String())

	return &GetWatchlistOutput{
		Watchlist: user.Watchlist,
	}, nil
}

package account

import (
	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/domain/valueobject"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type GetFavoritesInput struct {
	Email string
}

type GetFavoritesOutput struct {
	Favorites []models.Movie
}

type GetFavoritesUseCase struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewGetFavoritesUseCase(repo repository.UserRepository, log *logger.Logger) *GetFavoritesUseCase {
	return &GetFavoritesUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *GetFavoritesUseCase) Execute(input GetFavoritesInput) (*GetFavoritesOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepo.GetAccountDetails(email.String())
	if err != nil {
		uc.logger.Error("Failed to get user favorites", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved favorites for: " + email.String())

	return &GetFavoritesOutput{
		Favorites: user.Favorites,
	}, nil
}

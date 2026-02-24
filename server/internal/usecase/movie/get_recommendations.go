package movie

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type GetRecommendationsInput struct {
	Email string
}

type GetRecommendationsOutput struct {
	Movies []models.Movie
}

type GetRecommendationsUseCase struct {
	recRepo repository.RecommendationRepository
	logger  *logger.Logger
}

func NewGetRecommendationsUseCase(recRepo repository.RecommendationRepository, log *logger.Logger) *GetRecommendationsUseCase {
	return &GetRecommendationsUseCase{
		recRepo: recRepo,
		logger:  log,
	}
}

func (uc *GetRecommendationsUseCase) Execute(input GetRecommendationsInput) (*GetRecommendationsOutput, error) {
	userID, err := uc.recRepo.GetUserIDByEmail(input.Email)
	if err != nil {
		uc.logger.Error("Failed to get user ID for recommendations", err)
		return &GetRecommendationsOutput{Movies: []models.Movie{}}, nil
	}

	has, err := uc.recRepo.HasRecommendations(userID)
	if err != nil {
		uc.logger.Error("Failed to check recommendations", err)
		return &GetRecommendationsOutput{Movies: []models.Movie{}}, nil
	}

	if !has {
		return &GetRecommendationsOutput{Movies: []models.Movie{}}, nil
	}

	movies, err := uc.recRepo.GetRecommendations(userID, 20)
	if err != nil {
		uc.logger.Error("Failed to get recommendations", err)
		return &GetRecommendationsOutput{Movies: []models.Movie{}}, nil
	}

	uc.logger.Info("Successfully retrieved recommendations for user: " + input.Email)

	return &GetRecommendationsOutput{Movies: movies}, nil
}

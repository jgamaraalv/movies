package movie

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type UpdateUserRecommendationsInput struct {
	Email string
}

type UpdateUserRecommendationsOutput struct {
	Success bool
}

type UpdateUserRecommendationsUseCase struct {
	recRepo repository.RecommendationRepository
	logger  *logger.Logger
}

func NewUpdateUserRecommendationsUseCase(recRepo repository.RecommendationRepository, log *logger.Logger) *UpdateUserRecommendationsUseCase {
	return &UpdateUserRecommendationsUseCase{
		recRepo: recRepo,
		logger:  log,
	}
}

func (uc *UpdateUserRecommendationsUseCase) Execute(input UpdateUserRecommendationsInput) (*UpdateUserRecommendationsOutput, error) {
	userID, err := uc.recRepo.GetUserIDByEmail(input.Email)
	if err != nil {
		uc.logger.Error("Failed to get user ID for recommendation update", err)
		return &UpdateUserRecommendationsOutput{Success: false}, err
	}

	if err := uc.recRepo.RecomputeUserEmbedding(userID); err != nil {
		uc.logger.Error("Failed to recompute user embedding", err)
		return &UpdateUserRecommendationsOutput{Success: false}, err
	}

	if err := uc.recRepo.InvalidateRecommendations(userID); err != nil {
		uc.logger.Error("Failed to invalidate recommendations", err)
		return &UpdateUserRecommendationsOutput{Success: false}, err
	}

	uc.logger.Info("Successfully updated recommendations for user: " + input.Email)

	return &UpdateUserRecommendationsOutput{Success: true}, nil
}

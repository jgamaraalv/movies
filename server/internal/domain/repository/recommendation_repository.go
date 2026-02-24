package repository

import "github.com/jgamaraalv/movies.git/models"

type RecommendationRepository interface {
	GetRecommendations(userID int, limit int) ([]models.Movie, error)
	HasRecommendations(userID int) (bool, error)
	GetUserIDByEmail(email string) (int, error)
	InvalidateRecommendations(userID int) error
	RecomputeUserEmbedding(userID int) error
	ComputeRecommendations(userID int) error
}

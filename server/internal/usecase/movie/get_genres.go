package movie

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type GetGenresOutput struct {
	Genres []models.Genre
}

type GetGenresUseCase struct {
	movieRepo repository.MovieRepository
	logger    *logger.Logger
}

func NewGetGenresUseCase(repo repository.MovieRepository, log *logger.Logger) *GetGenresUseCase {
	return &GetGenresUseCase{
		movieRepo: repo,
		logger:    log,
	}
}

func (uc *GetGenresUseCase) Execute() (*GetGenresOutput, error) {
	genres, err := uc.movieRepo.GetAllGenres()
	if err != nil {
		uc.logger.Error("Failed to get genres", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved all genres")

	return &GetGenresOutput{
		Genres: genres,
	}, nil
}

package movie

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type GetTopMoviesOutput struct {
	Movies []models.Movie
}

type GetTopMoviesUseCase struct {
	movieRepo repository.MovieRepository
	logger    *logger.Logger
}

func NewGetTopMoviesUseCase(repo repository.MovieRepository, log *logger.Logger) *GetTopMoviesUseCase {
	return &GetTopMoviesUseCase{
		movieRepo: repo,
		logger:    log,
	}
}

func (uc *GetTopMoviesUseCase) Execute() (*GetTopMoviesOutput, error) {
	movies, err := uc.movieRepo.GetTopMovies()
	if err != nil {
		uc.logger.Error("Failed to get top movies", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved top movies")

	return &GetTopMoviesOutput{
		Movies: movies,
	}, nil
}

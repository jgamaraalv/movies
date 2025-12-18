package movie

import (
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type GetRandomMoviesOutput struct {
	Movies []models.Movie
}

type GetRandomMoviesUseCase struct {
	movieRepo repository.MovieRepository
	logger    *logger.Logger
}

func NewGetRandomMoviesUseCase(repo repository.MovieRepository, log *logger.Logger) *GetRandomMoviesUseCase {
	return &GetRandomMoviesUseCase{
		movieRepo: repo,
		logger:    log,
	}
}

func (uc *GetRandomMoviesUseCase) Execute() (*GetRandomMoviesOutput, error) {
	movies, err := uc.movieRepo.GetRandomMovies()
	if err != nil {
		uc.logger.Error("Failed to get random movies", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved random movies")

	return &GetRandomMoviesOutput{
		Movies: movies,
	}, nil
}

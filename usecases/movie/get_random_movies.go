package movie

import (
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
)

type GetRandomMoviesOutput struct {
	Movies []models.Movie
}

type GetRandomMoviesUseCase struct {
	movieStorage providers.MovieStorage
	logger       *logger.Logger
}

func NewGetRandomMoviesUseCase(storage providers.MovieStorage, log *logger.Logger) *GetRandomMoviesUseCase {
	return &GetRandomMoviesUseCase{
		movieStorage: storage,
		logger:       log,
	}
}

func (uc *GetRandomMoviesUseCase) Execute() (*GetRandomMoviesOutput, error) {
	movies, err := uc.movieStorage.GetRandomMovies()
	if err != nil {
		uc.logger.Error("Failed to get random movies", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved random movies")

	return &GetRandomMoviesOutput{
		Movies: movies,
	}, nil
}

package movie

import (
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
)

type GetTopMoviesOutput struct {
	Movies []models.Movie
}

type GetTopMoviesUseCase struct {
	movieStorage providers.MovieStorage
	logger       *logger.Logger
}

func NewGetTopMoviesUseCase(storage providers.MovieStorage, log *logger.Logger) *GetTopMoviesUseCase {
	return &GetTopMoviesUseCase{
		movieStorage: storage,
		logger:       log,
	}
}

func (uc *GetTopMoviesUseCase) Execute() (*GetTopMoviesOutput, error) {
	movies, err := uc.movieStorage.GetTopMovies()
	if err != nil {
		uc.logger.Error("Failed to get top movies", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved top movies")

	return &GetTopMoviesOutput{
		Movies: movies,
	}, nil
}

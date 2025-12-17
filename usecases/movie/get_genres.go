package movie

import (
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/providers"
)

type GetGenresOutput struct {
	Genres []models.Genre
}

type GetGenresUseCase struct {
	movieStorage providers.MovieStorage
	logger       *logger.Logger
}

func NewGetGenresUseCase(storage providers.MovieStorage, log *logger.Logger) *GetGenresUseCase {
	return &GetGenresUseCase{
		movieStorage: storage,
		logger:       log,
	}
}

func (uc *GetGenresUseCase) Execute() (*GetGenresOutput, error) {
	genres, err := uc.movieStorage.GetAllGenres()
	if err != nil {
		uc.logger.Error("Failed to get genres", err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved all genres")

	return &GetGenresOutput{
		Genres: genres,
	}, nil
}

package movie

import (
	"errors"

	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type SearchMoviesInput struct {
	Query string
	Order string
	Genre *int
}

type SearchMoviesOutput struct {
	Movies []models.Movie
}

type SearchMoviesUseCase struct {
	movieRepo repository.MovieRepository
	logger    *logger.Logger
}

func NewSearchMoviesUseCase(repo repository.MovieRepository, log *logger.Logger) *SearchMoviesUseCase {
	return &SearchMoviesUseCase{
		movieRepo: repo,
		logger:    log,
	}
}

func (uc *SearchMoviesUseCase) Execute(input SearchMoviesInput) (*SearchMoviesOutput, error) {
	if input.Query == "" {
		return nil, errors.New("search query is required")
	}

	movies, err := uc.movieRepo.SearchMoviesByName(input.Query, input.Order, input.Genre)
	if err != nil {
		uc.logger.Error("Failed to search movies", err)
		return nil, err
	}

	uc.logger.Info("Successfully searched movies with query: " + input.Query)

	return &SearchMoviesOutput{
		Movies: movies,
	}, nil
}

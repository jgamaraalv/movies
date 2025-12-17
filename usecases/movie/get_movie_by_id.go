package movie

import (
	"errors"
	"strconv"

	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type GetMovieByIDInput struct {
	ID int
}

type GetMovieByIDOutput struct {
	Movie models.Movie
}

type GetMovieByIDUseCase struct {
	movieRepo repository.MovieRepository
	logger    *logger.Logger
}

func NewGetMovieByIDUseCase(repo repository.MovieRepository, log *logger.Logger) *GetMovieByIDUseCase {
	return &GetMovieByIDUseCase{
		movieRepo: repo,
		logger:    log,
	}
}

func (uc *GetMovieByIDUseCase) Execute(input GetMovieByIDInput) (*GetMovieByIDOutput, error) {
	if input.ID <= 0 {
		return nil, errors.New("invalid movie ID")
	}

	movie, err := uc.movieRepo.GetMovieByID(input.ID)
	if err != nil {
		uc.logger.Error("Failed to get movie by ID: "+strconv.Itoa(input.ID), err)
		return nil, err
	}

	uc.logger.Info("Successfully retrieved movie with ID: " + strconv.Itoa(input.ID))

	return &GetMovieByIDOutput{
		Movie: movie,
	}, nil
}

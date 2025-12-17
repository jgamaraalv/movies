package movie

import (
	"errors"
	"strconv"

	"github.com/jgamaraalv/movies.git/domain/entity"
	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type GetMovieByIDInput struct {
	ID int
}

type MovieDetails struct {
	IsRecent       bool
	IsClassic      bool
	IsHighlyRated  bool
	IsPopular      bool
	HasTrailer     bool
	HasPoster      bool
	HasTagline     bool
	FormattedScore string
	TitleWithYear  string
	GenreNames     []string
	MainCast       []string
}

type GetMovieByIDOutput struct {
	Movie   models.Movie
	Details MovieDetails
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

	movieModel, err := uc.movieRepo.GetMovieByID(input.ID)
	if err != nil {
		uc.logger.Error("Failed to get movie by ID: "+strconv.Itoa(input.ID), err)
		return nil, err
	}

	movieEntity := entity.MovieFromModel(movieModel)

	mainCastNames := make([]string, 0)
	for _, actor := range movieEntity.MainCast(5) {
		mainCastNames = append(mainCastNames, actor.FullName())
	}

	details := MovieDetails{
		IsRecent:       movieEntity.IsRecent(),
		IsClassic:      movieEntity.IsClassic(),
		IsHighlyRated:  movieEntity.IsHighlyRated(),
		IsPopular:      movieEntity.IsPopular(),
		HasTrailer:     movieEntity.HasTrailer(),
		HasPoster:      movieEntity.HasPoster(),
		HasTagline:     movieEntity.HasTagline(),
		FormattedScore: movieEntity.FormattedScore(),
		TitleWithYear:  movieEntity.TitleWithYear(),
		GenreNames:     movieEntity.GenreNames(),
		MainCast:       mainCastNames,
	}

	uc.logger.Info("Successfully retrieved movie with ID: " + strconv.Itoa(input.ID))

	return &GetMovieByIDOutput{
		Movie:   movieModel,
		Details: details,
	}, nil
}

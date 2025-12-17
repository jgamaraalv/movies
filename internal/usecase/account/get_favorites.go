package account

import (
	"github.com/jgamaraalv/movies.git/internal/domain/entity"
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/internal/domain/valueobject"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type GetFavoritesInput struct {
	Email string
}

type FavoriteMovieInfo struct {
	Movie          models.Movie
	IsRecent       bool
	IsClassic      bool
	IsHighlyRated  bool
	HasTrailer     bool
	FormattedScore string
	TitleWithYear  string
}

type GetFavoritesOutput struct {
	Favorites        []models.Movie
	FavoritesInfo    []FavoriteMovieInfo
	TotalCount       int
	HighlyRatedCount int
	RecentCount      int
}

type GetFavoritesUseCase struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewGetFavoritesUseCase(repo repository.UserRepository, log *logger.Logger) *GetFavoritesUseCase {
	return &GetFavoritesUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *GetFavoritesUseCase) Execute(input GetFavoritesInput) (*GetFavoritesOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	userModel, err := uc.userRepo.GetAccountDetails(email.String())
	if err != nil {
		uc.logger.Error("Failed to get user favorites", err)
		return nil, err
	}

	favoritesInfo := make([]FavoriteMovieInfo, len(userModel.Favorites))
	highlyRatedCount := 0
	recentCount := 0

	for i, m := range userModel.Favorites {
		movieEntity := entity.MovieFromModel(m)

		favoritesInfo[i] = FavoriteMovieInfo{
			Movie:          m,
			IsRecent:       movieEntity.IsRecent(),
			IsClassic:      movieEntity.IsClassic(),
			IsHighlyRated:  movieEntity.IsHighlyRated(),
			HasTrailer:     movieEntity.HasTrailer(),
			FormattedScore: movieEntity.FormattedScore(),
			TitleWithYear:  movieEntity.TitleWithYear(),
		}

		if movieEntity.IsHighlyRated() {
			highlyRatedCount++
		}
		if movieEntity.IsRecent() {
			recentCount++
		}
	}

	uc.logger.Info("Successfully retrieved favorites for: " + email.String())

	return &GetFavoritesOutput{
		Favorites:        userModel.Favorites,
		FavoritesInfo:    favoritesInfo,
		TotalCount:       len(userModel.Favorites),
		HighlyRatedCount: highlyRatedCount,
		RecentCount:      recentCount,
	}, nil
}

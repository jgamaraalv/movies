package account

import (
	"github.com/jgamaraalv/movies.git/domain/entity"
	"github.com/jgamaraalv/movies.git/domain/repository"
	"github.com/jgamaraalv/movies.git/domain/valueobject"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/models"
)

type GetWatchlistInput struct {
	Email string
}

type WatchlistMovieInfo struct {
	Movie          models.Movie
	IsRecent       bool
	IsClassic      bool
	IsHighlyRated  bool
	HasTrailer     bool
	FormattedScore string
	TitleWithYear  string
}

type GetWatchlistOutput struct {
	Watchlist        []models.Movie
	WatchlistInfo    []WatchlistMovieInfo
	TotalCount       int
	HighlyRatedCount int
	RecentCount      int
}

type GetWatchlistUseCase struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewGetWatchlistUseCase(repo repository.UserRepository, log *logger.Logger) *GetWatchlistUseCase {
	return &GetWatchlistUseCase{
		userRepo: repo,
		logger:   log,
	}
}

func (uc *GetWatchlistUseCase) Execute(input GetWatchlistInput) (*GetWatchlistOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	userModel, err := uc.userRepo.GetAccountDetails(email.String())
	if err != nil {
		uc.logger.Error("Failed to get user watchlist", err)
		return nil, err
	}

	watchlistInfo := make([]WatchlistMovieInfo, len(userModel.Watchlist))
	highlyRatedCount := 0
	recentCount := 0

	for i, m := range userModel.Watchlist {
		movieEntity := entity.MovieFromModel(m)

		watchlistInfo[i] = WatchlistMovieInfo{
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

	uc.logger.Info("Successfully retrieved watchlist for: " + email.String())

	return &GetWatchlistOutput{
		Watchlist:        userModel.Watchlist,
		WatchlistInfo:    watchlistInfo,
		TotalCount:       len(userModel.Watchlist),
		HighlyRatedCount: highlyRatedCount,
		RecentCount:      recentCount,
	}, nil
}

package entity

import "github.com/jgamaraalv/movies.git/models"

func MovieFromModel(m models.Movie) *Movie {
	movie := ReconstructMovie(
		m.ID,
		m.TMDB_ID,
		m.Title,
		m.Tagline,
		m.ReleaseYear,
		m.Overview,
		m.Score,
		m.Popularity,
		m.Language,
		m.PosterURL,
		m.TrailerURL,
	)

	genres := make([]Genre, len(m.Genres))
	for i, g := range m.Genres {
		genres[i] = Genre{ID: g.ID, Name: g.Name}
	}
	movie.SetGenres(genres)

	casting := make([]Actor, len(m.Casting))
	for i, a := range m.Casting {
		casting[i] = Actor{
			ID:        a.ID,
			FirstName: a.FirstName,
			LastName:  a.LastName,
			ImageURL:  a.ImageURL,
		}
	}
	movie.SetCasting(casting)

	movie.SetKeywords(m.Keywords)

	return movie
}

func MovieToModel(m *Movie) models.Movie {
	genres := make([]models.Genre, len(m.Genres()))
	for i, g := range m.Genres() {
		genres[i] = models.Genre{ID: g.ID, Name: g.Name}
	}

	casting := make([]models.Actor, len(m.Casting()))
	for i, a := range m.Casting() {
		casting[i] = models.Actor{
			ID:        a.ID,
			FirstName: a.FirstName,
			LastName:  a.LastName,
			ImageURL:  a.ImageURL,
		}
	}

	return models.Movie{
		ID:          m.ID(),
		TMDB_ID:     m.TMDBID(),
		Title:       m.Title(),
		Tagline:     m.Tagline(),
		ReleaseYear: m.ReleaseYear(),
		Genres:      genres,
		Overview:    m.Overview(),
		Score:       m.Score(),
		Popularity:  m.Popularity(),
		Keywords:    m.Keywords(),
		Language:    m.Language(),
		PosterURL:   m.PosterURL(),
		TrailerURL:  m.TrailerURL(),
		Casting:     casting,
	}
}

func MoviesFromModels(models []models.Movie) []*Movie {
	entities := make([]*Movie, len(models))
	for i, m := range models {
		entities[i] = MovieFromModel(m)
	}
	return entities
}

func UserFromModel(m models.User) (*User, error) {
	favoriteIDs := make([]int, len(m.Favorites))
	for i, movie := range m.Favorites {
		favoriteIDs[i] = movie.ID
	}

	watchlistIDs := make([]int, len(m.Watchlist))
	for i, movie := range m.Watchlist {
		watchlistIDs[i] = movie.ID
	}

	return ReconstructUser(m.ID, m.Name, m.Email, favoriteIDs, watchlistIDs)
}

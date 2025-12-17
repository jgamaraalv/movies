package entity

import (
	"fmt"
	"time"
)

type Genre struct {
	ID   int
	Name string
}

type Actor struct {
	ID        int
	FirstName string
	LastName  string
	ImageURL  *string
}

func (a Actor) FullName() string {
	return a.FirstName + " " + a.LastName
}

func (a Actor) HasImage() bool {
	return a.ImageURL != nil && *a.ImageURL != ""
}

type Movie struct {
	id          int
	tmdbID      int
	title       string
	tagline     *string
	releaseYear int
	genres      []Genre
	overview    *string
	score       *float32
	popularity  *float32
	keywords    []string
	language    *string
	posterURL   *string
	trailerURL  *string
	casting     []Actor
}

func NewMovie(id int, title string, releaseYear int) *Movie {
	return &Movie{
		id:          id,
		title:       title,
		releaseYear: releaseYear,
		genres:      make([]Genre, 0),
		keywords:    make([]string, 0),
		casting:     make([]Actor, 0),
	}
}

func ReconstructMovie(
	id, tmdbID int,
	title string,
	tagline *string,
	releaseYear int,
	overview *string,
	score, popularity *float32,
	language, posterURL, trailerURL *string,
) *Movie {
	return &Movie{
		id:          id,
		tmdbID:      tmdbID,
		title:       title,
		tagline:     tagline,
		releaseYear: releaseYear,
		overview:    overview,
		score:       score,
		popularity:  popularity,
		language:    language,
		posterURL:   posterURL,
		trailerURL:  trailerURL,
		genres:      make([]Genre, 0),
		keywords:    make([]string, 0),
		casting:     make([]Actor, 0),
	}
}

// Getters
func (m *Movie) ID() int              { return m.id }
func (m *Movie) TMDBID() int          { return m.tmdbID }
func (m *Movie) Title() string        { return m.title }
func (m *Movie) Tagline() *string     { return m.tagline }
func (m *Movie) ReleaseYear() int     { return m.releaseYear }
func (m *Movie) Overview() *string    { return m.overview }
func (m *Movie) Score() *float32      { return m.score }
func (m *Movie) Popularity() *float32 { return m.popularity }
func (m *Movie) Language() *string    { return m.language }
func (m *Movie) PosterURL() *string   { return m.posterURL }
func (m *Movie) TrailerURL() *string  { return m.trailerURL }

func (m *Movie) Genres() []Genre {
	result := make([]Genre, len(m.genres))
	copy(result, m.genres)
	return result
}

func (m *Movie) Keywords() []string {
	result := make([]string, len(m.keywords))
	copy(result, m.keywords)
	return result
}

func (m *Movie) Casting() []Actor {
	result := make([]Actor, len(m.casting))
	copy(result, m.casting)
	return result
}

// Setters for reconstruction
func (m *Movie) SetGenres(genres []Genre)      { m.genres = genres }
func (m *Movie) SetKeywords(keywords []string) { m.keywords = keywords }
func (m *Movie) SetCasting(casting []Actor)    { m.casting = casting }

// Behavior methods
func (m *Movie) HasTrailer() bool {
	return m.trailerURL != nil && *m.trailerURL != ""
}

func (m *Movie) HasPoster() bool {
	return m.posterURL != nil && *m.posterURL != ""
}

func (m *Movie) HasTagline() bool {
	return m.tagline != nil && *m.tagline != ""
}

func (m *Movie) HasOverview() bool {
	return m.overview != nil && *m.overview != ""
}

func (m *Movie) IsRecent() bool {
	currentYear := time.Now().Year()
	return m.releaseYear >= currentYear-2
}

func (m *Movie) IsClassic() bool {
	currentYear := time.Now().Year()
	return m.releaseYear <= currentYear-25
}

func (m *Movie) IsHighlyRated() bool {
	return m.score != nil && *m.score >= 7.5
}

func (m *Movie) IsPopular() bool {
	return m.popularity != nil && *m.popularity >= 100
}

func (m *Movie) GenreNames() []string {
	names := make([]string, len(m.genres))
	for i, g := range m.genres {
		names[i] = g.Name
	}
	return names
}

func (m *Movie) HasGenre(genreName string) bool {
	for _, g := range m.genres {
		if g.Name == genreName {
			return true
		}
	}
	return false
}

func (m *Movie) CastingNames() []string {
	names := make([]string, len(m.casting))
	for i, a := range m.casting {
		names[i] = a.FullName()
	}
	return names
}

func (m *Movie) MainCast(limit int) []Actor {
	if limit >= len(m.casting) {
		return m.Casting()
	}
	result := make([]Actor, limit)
	copy(result, m.casting[:limit])
	return result
}

func (m *Movie) FormattedScore() string {
	if m.score == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.1f/10", *m.score)
}

func (m *Movie) YearString() string {
	return fmt.Sprintf("%d", m.releaseYear)
}

func (m *Movie) TitleWithYear() string {
	return fmt.Sprintf("%s (%d)", m.title, m.releaseYear)
}

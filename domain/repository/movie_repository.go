package repository

import "github.com/jgamaraalv/movies.git/models"

type MovieRepository interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	SearchMoviesByName(query string, orderBy string, genreID *int) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}

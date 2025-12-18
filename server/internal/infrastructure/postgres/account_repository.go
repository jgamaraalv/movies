package postgres

import (
	"database/sql"
	"time"

	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAccountRepository(db *sql.DB, log *logger.Logger) (*AccountRepository, error) {
	return &AccountRepository{
		db:     db,
		logger: log,
	}, nil
}

func (r *AccountRepository) Register(name, email, hashedPassword string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing user", err)
		return false, err
	}
	if exists {
		r.logger.Error("User already exists with email: "+email, repository.ErrUserAlreadyExists)
		return false, repository.ErrUserAlreadyExists
	}

	query := `
		INSERT INTO users (name, email, password_hashed, time_created)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var userID int
	err = r.db.QueryRow(
		query,
		name,
		email,
		hashedPassword,
		time.Now(),
	).Scan(&userID)
	if err != nil {
		r.logger.Error("Failed to register user", err)
		return false, err
	}

	return true, nil
}

func (r *AccountRepository) Authenticate(email string, password string) (bool, error) {
	var user models.User
	query := `
		SELECT id, name, email, password_hashed
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return false, repository.ErrAuthenticationValidation
	}
	if err != nil {
		r.logger.Error("Failed to query user for authentication", err)
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		r.logger.Error("Password mismatch for email: "+email, nil)
		return false, repository.ErrAuthenticationValidation
	}

	updateQuery := `
		UPDATE users 
		SET last_login = $1
		WHERE id = $2
	`
	_, err = r.db.Exec(updateQuery, time.Now(), user.ID)
	if err != nil {
		r.logger.Error("Failed to update last login", err)
	}

	return true, nil
}

func (r *AccountRepository) GetAccountDetails(email string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return models.User{}, repository.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user by email", err)
		return models.User{}, err
	}

	favoritesQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'favorite'
	`
	favoriteRows, err := r.db.Query(favoritesQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user favorites", err)
		return user, err
	}
	defer favoriteRows.Close()

	for favoriteRows.Next() {
		var m models.Movie
		if err := favoriteRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan favorite movie row", err)
			return user, err
		}
		user.Favorites = append(user.Favorites, m)
	}

	watchlistQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'watchlist'
	`
	watchlistRows, err := r.db.Query(watchlistQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user watchlist", err)
		return user, err
	}
	defer watchlistRows.Close()

	for watchlistRows.Next() {
		var m models.Movie
		if err := watchlistRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan watchlist movie row", err)
			return user, err
		}
		user.Watchlist = append(user.Watchlist, m)
	}

	return user, nil
}

func (r *AccountRepository) SaveCollection(user models.User, movieID int, collection string) (bool, error) {
	var userID int
	err := r.db.QueryRow(`
		SELECT id 
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`, user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found", nil)
		return false, repository.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user ID", err)
		return false, err
	}

	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM user_movies 
			WHERE user_id = $1 
			AND movie_id = $2 
			AND relation_type = $3
		)
	`, userID, movieID, collection).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing collection entry", err)
		return false, err
	}
	if exists {
		r.logger.Info("Movie already in " + collection + " for user")
		return true, nil
	}

	query := `
		INSERT INTO user_movies (user_id, movie_id, relation_type, time_added)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.Exec(query, userID, movieID, collection, time.Now())
	if err != nil {
		r.logger.Error("Failed to save movie to "+collection, err)
		return false, err
	}

	return true, nil
}

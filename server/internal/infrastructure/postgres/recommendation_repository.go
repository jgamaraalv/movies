package postgres

import (
	"database/sql"

	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/models"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type RecommendationRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewRecommendationRepository(db *sql.DB, log *logger.Logger) (*RecommendationRepository, error) {
	return &RecommendationRepository{
		db:     db,
		logger: log,
	}, nil
}

func (r *RecommendationRepository) GetRecommendations(userID int, limit int) ([]models.Movie, error) {
	query := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year,
		       m.overview, m.score, m.popularity, m.language,
		       m.poster_url, m.trailer_url
		FROM user_recommendations ur
		JOIN movies m ON m.id = ur.movie_id
		WHERE ur.user_id = $1
		ORDER BY ur.score DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		r.logger.Error("Failed to query recommendations", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan recommendation movie row", err)
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func (r *RecommendationRepository) HasRecommendations(userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM user_recommendations WHERE user_id = $1)
	`, userID).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check recommendations existence", err)
		return false, err
	}
	return exists, nil
}

func (r *RecommendationRepository) GetUserIDByEmail(email string) (int, error) {
	var userID int
	err := r.db.QueryRow(`
		SELECT id FROM users WHERE email = $1 AND time_deleted IS NULL
	`, email).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, repository.ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to get user ID by email", err)
		return 0, err
	}
	return userID, nil
}

func (r *RecommendationRepository) InvalidateRecommendations(userID int) error {
	_, err := r.db.Exec(`DELETE FROM user_recommendations WHERE user_id = $1`, userID)
	if err != nil {
		r.logger.Error("Failed to invalidate recommendations", err)
		return err
	}
	return nil
}

func (r *RecommendationRepository) RecomputeUserEmbedding(userID int) error {
	// Average of movie embeddings for user's collections using pgvector's native AVG
	query := `
		INSERT INTO user_embeddings (user_id, embedding, updated_at)
		SELECT
			$1,
			avg(me.embedding)::vector(128),
			CURRENT_TIMESTAMP
		FROM user_movies um
		JOIN movie_embeddings me ON me.movie_id = um.movie_id
		WHERE um.user_id = $1
		HAVING avg(me.embedding) IS NOT NULL
		ON CONFLICT (user_id)
		DO UPDATE SET
			embedding = EXCLUDED.embedding,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		r.logger.Error("Failed to recompute user embedding", err)
		return err
	}
	return nil
}

package postgres

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"

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

type movieCandidate struct {
	movieID       int
	score         float64
	popularity    float64
	genreIDs      []int
	embSimilarity float64
	collabScore   float64
	genreAffinity float64
	finalScore    float64
}

func (r *RecommendationRepository) ComputeRecommendations(userID int) error {
	// 1. Delete existing recommendations
	if _, err := r.db.Exec(`DELETE FROM user_recommendations WHERE user_id = $1`, userID); err != nil {
		r.logger.Error("Failed to clear old recommendations", err)
		return err
	}

	// 2. Get user's movies to exclude
	rows, err := r.db.Query(`SELECT movie_id FROM user_movies WHERE user_id = $1`, userID)
	if err != nil {
		r.logger.Error("Failed to get user movies", err)
		return err
	}
	userMovieIDs := make(map[int]bool)
	for rows.Next() {
		var mid int
		rows.Scan(&mid)
		userMovieIDs[mid] = true
	}
	rows.Close()

	if len(userMovieIDs) == 0 {
		return nil
	}

	// 3. Get user's genre preferences (weighted by frequency)
	genreRows, err := r.db.Query(`
		SELECT mg.genre_id, COUNT(*) as cnt
		FROM user_movies um
		JOIN movie_genres mg ON mg.movie_id = um.movie_id
		WHERE um.user_id = $1
		GROUP BY mg.genre_id
		ORDER BY cnt DESC
	`, userID)
	if err != nil {
		r.logger.Error("Failed to get user genre preferences", err)
		return err
	}
	genreWeights := make(map[int]float64)
	var totalGenreCount float64
	for genreRows.Next() {
		var genreID, cnt int
		genreRows.Scan(&genreID, &cnt)
		genreWeights[genreID] = float64(cnt)
		totalGenreCount += float64(cnt)
	}
	genreRows.Close()

	// Normalize genre weights to [0, 1]
	if totalGenreCount > 0 {
		for gid := range genreWeights {
			genreWeights[gid] = genreWeights[gid] / totalGenreCount
		}
	}

	candidates := make(map[int]*movieCandidate)

	// 4. Embedding-based candidates (content-based via pgvector cosine distance)
	var hasEmbedding bool
	err = r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_embeddings WHERE user_id = $1)`, userID).Scan(&hasEmbedding)
	if err != nil {
		r.logger.Error("Failed to check user embedding", err)
		return err
	}

	if hasEmbedding {
		embRows, err := r.db.Query(`
			SELECT me.movie_id, 1 - (me.embedding <=> ue.embedding) as similarity
			FROM movie_embeddings me
			CROSS JOIN user_embeddings ue
			WHERE ue.user_id = $1
			AND me.movie_id NOT IN (SELECT movie_id FROM user_movies WHERE user_id = $1)
			ORDER BY me.embedding <=> ue.embedding
			LIMIT 50
		`, userID)
		if err != nil {
			r.logger.Error("Failed to get embedding candidates", err)
		} else {
			for embRows.Next() {
				var mid int
				var sim float64
				embRows.Scan(&mid, &sim)
				candidates[mid] = &movieCandidate{movieID: mid, embSimilarity: sim}
			}
			embRows.Close()
		}

		// 5. Collaborative filtering: similar users' movies
		collabRows, err := r.db.Query(`
			SELECT um.movie_id, um.relation_type
			FROM user_embeddings ue_other
			JOIN user_movies um ON um.user_id = ue_other.user_id
			CROSS JOIN user_embeddings ue_self
			WHERE ue_self.user_id = $1
			AND ue_other.user_id != $1
			AND um.movie_id NOT IN (SELECT movie_id FROM user_movies WHERE user_id = $1)
			ORDER BY ue_other.embedding <=> ue_self.embedding
			LIMIT 200
		`, userID)
		if err != nil {
			r.logger.Error("Failed to get collaborative candidates", err)
		} else {
			for collabRows.Next() {
				var mid int
				var relType string
				collabRows.Scan(&mid, &relType)
				weight := 1.0
				if relType == "watchlist" {
					weight = 0.5
				}
				if c, exists := candidates[mid]; exists {
					c.collabScore += weight
				} else {
					candidates[mid] = &movieCandidate{movieID: mid, collabScore: weight}
				}
			}
			collabRows.Close()
		}
	}

	// 6. Genre-based candidates: find movies matching user's preferred genres
	// This works even without embeddings and accesses the full movie catalog
	if len(genreWeights) > 0 {
		genreIDs := make([]string, 0, len(genreWeights))
		for gid := range genreWeights {
			genreIDs = append(genreIDs, fmt.Sprintf("%d", gid))
		}

		genreQuery := fmt.Sprintf(`
			SELECT m.id, m.score, m.popularity,
				   array_agg(mg.genre_id) as genre_ids
			FROM movies m
			JOIN movie_genres mg ON mg.movie_id = m.id
			WHERE mg.genre_id IN (%s)
			AND m.id NOT IN (SELECT movie_id FROM user_movies WHERE user_id = $1)
			GROUP BY m.id, m.score, m.popularity
			ORDER BY COUNT(DISTINCT mg.genre_id) DESC, m.score DESC
			LIMIT 200
		`, strings.Join(genreIDs, ","))

		genreCandRows, err := r.db.Query(genreQuery, userID)
		if err != nil {
			r.logger.Error("Failed to get genre candidates", err)
		} else {
			for genreCandRows.Next() {
				var mid int
				var movieScore, moviePop sql.NullFloat64
				var genreIDsStr string
				if err := genreCandRows.Scan(&mid, &movieScore, &moviePop, &genreIDsStr); err != nil {
					continue
				}
				if c, exists := candidates[mid]; exists {
					if movieScore.Valid {
						c.score = movieScore.Float64
					}
					if moviePop.Valid {
						c.popularity = moviePop.Float64
					}
					c.genreIDs = parseIntArray(genreIDsStr)
				} else {
					c := &movieCandidate{movieID: mid}
					if movieScore.Valid {
						c.score = movieScore.Float64
					}
					if moviePop.Valid {
						c.popularity = moviePop.Float64
					}
					c.genreIDs = parseIntArray(genreIDsStr)
					candidates[mid] = c
				}
			}
			genreCandRows.Close()
		}
	}

	// 7. For candidates missing score/popularity/genres, fetch them
	missingIDs := make([]int, 0)
	for mid, c := range candidates {
		if len(c.genreIDs) == 0 {
			missingIDs = append(missingIDs, mid)
		}
	}
	if len(missingIDs) > 0 {
		placeholders := make([]string, len(missingIDs))
		args := make([]interface{}, len(missingIDs))
		for i, mid := range missingIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = mid
		}
		metaRows, err := r.db.Query(fmt.Sprintf(`
			SELECT m.id, m.score, m.popularity, COALESCE(
				(SELECT array_agg(mg.genre_id) FROM movie_genres mg WHERE mg.movie_id = m.id),
				'{}'
			)
			FROM movies m WHERE m.id IN (%s)
		`, strings.Join(placeholders, ",")), args...)
		if err == nil {
			for metaRows.Next() {
				var mid int
				var movieScore, moviePop sql.NullFloat64
				var genreIDsStr string
				if metaRows.Scan(&mid, &movieScore, &moviePop, &genreIDsStr) == nil {
					if c, exists := candidates[mid]; exists {
						if movieScore.Valid {
							c.score = movieScore.Float64
						}
						if moviePop.Valid {
							c.popularity = moviePop.Float64
						}
						c.genreIDs = parseIntArray(genreIDsStr)
					}
				}
			}
			metaRows.Close()
		}
	}

	// 8. Compute genre affinity for each candidate
	for _, c := range candidates {
		if len(c.genreIDs) > 0 && len(genreWeights) > 0 {
			var affinity float64
			for _, gid := range c.genreIDs {
				if w, ok := genreWeights[gid]; ok {
					affinity += w
				}
			}
			c.genreAffinity = affinity
		}
	}

	// 9. Compute final blended score
	// Normalize collaborative scores
	var maxCollab float64
	for _, c := range candidates {
		if c.collabScore > maxCollab {
			maxCollab = c.collabScore
		}
	}

	for _, c := range candidates {
		normCollab := 0.0
		if maxCollab > 0 {
			normCollab = c.collabScore / maxCollab
		}
		normEmb := math.Max(0, c.embSimilarity)
		normScore := c.score / 10.0
		normPop := math.Min(c.popularity, 100.0) / 100.0

		// Genre affinity is the strongest signal for taste matching
		// Embedding similarity captures latent features from the NCF model
		// Collaborative filtering captures what similar users liked
		// Movie score and popularity are tiebreakers
		c.finalScore = c.genreAffinity*0.35 +
			normEmb*0.25 +
			normCollab*0.20 +
			normScore*0.12 +
			normPop*0.08
	}

	// 10. Sort and pick top 20
	candidateList := make([]*movieCandidate, 0, len(candidates))
	for _, c := range candidates {
		candidateList = append(candidateList, c)
	}
	sort.Slice(candidateList, func(i, j int) bool {
		return candidateList[i].finalScore > candidateList[j].finalScore
	})

	limit := 20
	if len(candidateList) < limit {
		limit = len(candidateList)
	}

	// 11. Store recommendations
	for _, c := range candidateList[:limit] {
		_, err := r.db.Exec(`
			INSERT INTO user_recommendations (user_id, movie_id, score, computed_at)
			VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
			ON CONFLICT (user_id, movie_id)
			DO UPDATE SET score = EXCLUDED.score, computed_at = CURRENT_TIMESTAMP
		`, userID, c.movieID, c.finalScore)
		if err != nil {
			r.logger.Error("Failed to insert recommendation", err)
		}
	}

	r.logger.Info(fmt.Sprintf("Computed %d recommendations for user %d (from %d candidates)", limit, userID, len(candidates)))
	return nil
}

// parseIntArray parses a PostgreSQL int array string like "{1,2,3}" into a slice of ints
func parseIntArray(s string) []int {
	s = strings.Trim(s, "{}")
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		var v int
		if _, err := fmt.Sscanf(strings.TrimSpace(p), "%d", &v); err == nil {
			result = append(result, v)
		}
	}
	return result
}

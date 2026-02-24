package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	"github.com/jgamaraalv/movies.git/internal/usecase/movie"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

type MovieHandler struct {
	getTopMoviesUC       *movie.GetTopMoviesUseCase
	getRandomMoviesUC    *movie.GetRandomMoviesUseCase
	searchMoviesUC       *movie.SearchMoviesUseCase
	getMovieByIDUC       *movie.GetMovieByIDUseCase
	getGenresUC          *movie.GetGenresUseCase
	getRecommendationsUC *movie.GetRecommendationsUseCase
	logger               *logger.Logger
}

func NewMovieHandler(repo repository.MovieRepository, recRepo repository.RecommendationRepository, log *logger.Logger) *MovieHandler {
	h := &MovieHandler{
		getTopMoviesUC:    movie.NewGetTopMoviesUseCase(repo, log),
		getRandomMoviesUC: movie.NewGetRandomMoviesUseCase(repo, log),
		searchMoviesUC:    movie.NewSearchMoviesUseCase(repo, log),
		getMovieByIDUC:    movie.NewGetMovieByIDUseCase(repo, log),
		getGenresUC:       movie.NewGetGenresUseCase(repo, log),
		logger:            log,
	}
	if recRepo != nil {
		h.getRecommendationsUC = movie.NewGetRecommendationsUseCase(recRepo, log)
	}
	return h
}

func (h *MovieHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *MovieHandler) handleError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		if err == repository.ErrMovieNotFound {
			http.Error(w, context, http.StatusNotFound)
			return true
		}
		h.logger.Error(context, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

func (h *MovieHandler) parseID(w http.ResponseWriter, idStr string) (int, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid ID format", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	output, err := h.getTopMoviesUC.Execute()
	if h.handleError(w, err, "Failed to get top movies") {
		return
	}
	h.writeJSONResponse(w, output.Movies)
}

func (h *MovieHandler) GetRandomMovies(w http.ResponseWriter, r *http.Request) {
	output, err := h.getRandomMoviesUC.Execute()
	if h.handleError(w, err, "Failed to get random movies") {
		return
	}
	h.writeJSONResponse(w, output.Movies)
}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")
	genreStr := r.URL.Query().Get("genre")

	var genre *int
	if genreStr != "" {
		genreInt, ok := h.parseID(w, genreStr)
		if !ok {
			return
		}
		genre = &genreInt
	}

	input := movie.SearchMoviesInput{
		Query: query,
		Order: order,
		Genre: genre,
	}

	output, err := h.searchMoviesUC.Execute(input)
	if err != nil {
		h.writeJSONResponse(w, []interface{}{})
		return
	}
	h.writeJSONResponse(w, output.Movies)
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/movies/"):]
	id, ok := h.parseID(w, idStr)
	if !ok {
		return
	}

	input := movie.GetMovieByIDInput{ID: id}
	output, err := h.getMovieByIDUC.Execute(input)
	if h.handleError(w, err, "Failed to get movie by ID") {
		return
	}
	h.writeJSONResponse(w, output.Movie)
}

func (h *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	output, err := h.getGenresUC.Execute()
	if h.handleError(w, err, "Failed to get genres") {
		return
	}
	h.writeJSONResponse(w, output.Genres)
}

func (h *MovieHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	if h.getRecommendationsUC == nil {
		h.writeJSONResponse(w, []interface{}{})
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		h.writeJSONResponse(w, []interface{}{})
		return
	}

	input := movie.GetRecommendationsInput{Email: email}
	output, err := h.getRecommendationsUC.Execute(input)
	if err != nil {
		h.writeJSONResponse(w, []interface{}{})
		return
	}
	h.writeJSONResponse(w, output.Movies)
}

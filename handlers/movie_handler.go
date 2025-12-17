package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/providers"
)

type MovieHandler struct {
	storage providers.MovieStorage
	logger  *logger.Logger
}

func (h *MovieHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		if err == providers.ErrMovieNotFound {
			http.Error(w, context, http.StatusNotFound)
			return true
		}
		h.logger.Error(context, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
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

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.storage.GetTopMovies()
	if h.handleStorageError(w, err, "Failed to get top movies") {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Error("Failed to get top movies", err)
	}

	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served top movies")
	}
}

func NewMovieHandler(storage providers.MovieStorage, log *logger.Logger) *MovieHandler {
	return &MovieHandler{
		storage: storage,
		logger: log,
	}
}

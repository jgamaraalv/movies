package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jgamaraalv/movies.git/internal/domain/repository"
	accountuc "github.com/jgamaraalv/movies.git/internal/usecase/account"
	movieuc "github.com/jgamaraalv/movies.git/internal/usecase/movie"
	"github.com/jgamaraalv/movies.git/pkg/logger"
	"github.com/jgamaraalv/movies.git/pkg/token"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CollectionRequest struct {
	MovieID    int    `json:"movie_id"`
	Collection string `json:"collection"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	JWT     string `json:"jwt,omitempty"`
}

type AccountHandler struct {
	registerUC                   *accountuc.RegisterUseCase
	authenticateUC               *accountuc.AuthenticateUseCase
	getFavoritesUC               *accountuc.GetFavoritesUseCase
	getWatchlistUC               *accountuc.GetWatchlistUseCase
	saveToCollectionUC           *accountuc.SaveToCollectionUseCase
	updateUserRecommendationsUC  *movieuc.UpdateUserRecommendationsUseCase
	logger                       *logger.Logger
}

func NewAccountHandler(repo repository.UserRepository, recRepo repository.RecommendationRepository, log *logger.Logger) *AccountHandler {
	h := &AccountHandler{
		registerUC:         accountuc.NewRegisterUseCase(repo, log),
		authenticateUC:     accountuc.NewAuthenticateUseCase(repo, log),
		getFavoritesUC:     accountuc.NewGetFavoritesUseCase(repo, log),
		getWatchlistUC:     accountuc.NewGetWatchlistUseCase(repo, log),
		saveToCollectionUC: accountuc.NewSaveToCollectionUseCase(repo, log),
		logger:             log,
	}
	if recRepo != nil {
		h.updateUserRecommendationsUC = movieuc.NewUpdateUserRecommendationsUseCase(recRepo, log)
	}
	return h
}

func (h *AccountHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *AccountHandler) handleError(w http.ResponseWriter, err error) bool {
	if err != nil {
		switch err {
		case repository.ErrAuthenticationValidation, repository.ErrUserAlreadyExists, repository.ErrRegistrationValidation:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return true
		case repository.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
			return true
		default:
			h.logger.Error("Handler error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return true
		}
	}
	return false
}

func (h *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode registration request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := accountuc.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.registerUC.Execute(input)
	if h.handleError(w, err) {
		return
	}

	response := AuthResponse{
		Success: output.Success,
		Message: output.Message,
		JWT:     output.JWT,
	}
	h.writeJSONResponse(w, response)
}

func (h *AccountHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode authentication request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := accountuc.AuthenticateInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.authenticateUC.Execute(input)
	if h.handleError(w, err) {
		return
	}

	response := AuthResponse{
		Success: output.Success,
		Message: output.Message,
		JWT:     output.JWT,
	}
	h.writeJSONResponse(w, response)
}

func (h *AccountHandler) SaveToCollection(w http.ResponseWriter, r *http.Request) {
	var req CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	input := accountuc.SaveToCollectionInput{
		Email:      email,
		MovieID:    req.MovieID,
		Collection: req.Collection,
	}

	output, err := h.saveToCollectionUC.Execute(input)
	if h.handleError(w, err) {
		return
	}

	// Fire-and-forget: recompute user embedding and invalidate recommendations
	if h.updateUserRecommendationsUC != nil && !output.AlreadyInCollection {
		go func() {
			recInput := movieuc.UpdateUserRecommendationsInput{Email: email}
			if _, err := h.updateUserRecommendationsUC.Execute(recInput); err != nil {
				h.logger.Error("Failed to update user recommendations in background", err)
			}
		}()
	}

	response := AuthResponse{
		Success: output.Success,
		Message: output.Message,
	}
	h.writeJSONResponse(w, response)
}

func (h *AccountHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	input := accountuc.GetFavoritesInput{Email: email}
	output, err := h.getFavoritesUC.Execute(input)
	if h.handleError(w, err) {
		return
	}

	h.writeJSONResponse(w, output.Favorites)
}

func (h *AccountHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	input := accountuc.GetWatchlistInput{Email: email}
	output, err := h.getWatchlistUC.Execute(input)
	if h.handleError(w, err) {
		return
	}

	h.writeJSONResponse(w, output.Watchlist)
}

func (h *AccountHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		parsedToken, err := jwt.Parse(tokenStr,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(token.GetJWTSecret(*h.logger)), nil
			},
		)
		if err != nil || !parsedToken.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Email not found in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

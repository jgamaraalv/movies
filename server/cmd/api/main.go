package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jgamaraalv/movies.git/internal/handler"
	"github.com/jgamaraalv/movies.git/internal/infrastructure/postgres"
	"github.com/jgamaraalv/movies.git/pkg/logger"
)

func main() {
	logInstance := initializeLogger()

	// Try to load .env from multiple locations
	// Priority: current dir (.env) -> parent dir (../.env) -> grandparent dir (../../.env)
	envLoaded := false
	envPaths := []string{".env", "../.env", "../../.env"}

	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			log.Printf("Loaded .env from: %s", path)
			break
		}
	}

	if !envLoaded {
		log.Printf("No .env file found in any of the expected locations. Using environment variables only.")
	}

	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatalf("DATABASE_URL not set in environment")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	movieRepo, err := postgres.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository: %v", err)
	}

	accountRepo, err := postgres.NewAccountRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize account repository: %v", err)
	}

	// Initialize handlers
	movieHandler := handler.NewMovieHandler(movieRepo, logInstance)
	accountHandler := handler.NewAccountHandler(accountRepo, logInstance)

	// Health check endpoint for Docker/Kubernetes probes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","error":"database connection failed"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Set up routes
	http.HandleFunc("/api/account/register/", accountHandler.Register)
	http.HandleFunc("/api/account/authenticate/", accountHandler.Authenticate)
	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/random", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/search", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie)
	http.HandleFunc("/api/genres", movieHandler.GetGenres)

	http.Handle("/api/account/favorites/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetFavorites)))

	http.Handle("/api/account/watchlist/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetWatchlist)))

	http.Handle("/api/account/save-to-collection/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.SaveToCollection)))

	// Get public directory path (from root of project)
	publicDir := os.Getenv("PUBLIC_DIR")
	if publicDir == "" {
		// Default: assume we're running from project root, or use relative path
		publicDir = "public"
		// Try relative path from server/cmd/api (for development)
		if _, err := os.Stat(publicDir); os.IsNotExist(err) {
			publicDir = "../../public"
		}
	}
	publicDir, _ = filepath.Abs(publicDir)

	catchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
	}
	http.HandleFunc("/movies", catchAllHandler)
	http.HandleFunc("/movies/", catchAllHandler)
	http.HandleFunc("/account/", catchAllHandler)

	http.Handle("/", http.FileServer(http.Dir(publicDir)))

	const addr = ":8080"
	logInstance.Info("Server starting on " + addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logInstance.Error("Server failed to start", err)
		log.Fatalf("Server failed: %v", err)
	}
}

func initializeLogger() *logger.Logger {
	// Define o caminho do arquivo de log
	// Em produção com filesystem read-only, usa /app/logs ou /tmp
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		// Tenta /app/logs primeiro (tmpfs em produção)
		if _, err := os.Stat("/app/logs"); err == nil {
			logPath = "/app/logs/movie-service.log"
		} else if _, err := os.Stat("/tmp"); err == nil {
			// Fallback para /tmp
			logPath = "/tmp/movie-service.log"
		} else {
			// Último fallback: diretório atual (desenvolvimento)
			logPath = "movie-service.log"
		}
	}

	logInstance, err := logger.NewLogger(logPath)
	if err != nil {
		log.Printf("Warning: Failed to initialize file logger at %s: %v. Using stdout only.", logPath, err)
		// Fallback: cria um logger que só usa stdout
		logInstance, _ = logger.NewLogger("/dev/null")
	}
	return logInstance
}

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jgamaraalv/movies.git/handlers"
	"github.com/jgamaraalv/movies.git/logger"
	"github.com/jgamaraalv/movies.git/providers"
)

func main() {
	// Log Initializer
	logInstance := initializeLogger()

	// Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
	}

	// Database connection
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
	movieRepo, err := providers.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository: %v", err)
	}

	// Movie Handler Initializer
	movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)

	// Set up routes
	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/random", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/search", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie)
	http.HandleFunc("/api/genres", movieHandler.GetGenres)

	// Handler for client SPA
	catchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}
	http.HandleFunc("/movies", catchAllHandler)
	http.HandleFunc("/movies/", catchAllHandler)
	http.HandleFunc("/account/", catchAllHandler)

	// Handler for static files (frontend)
	http.Handle("/", http.FileServer(http.Dir("public")))

	const addr = ":8080"
	logInstance.Info("Server starting on " + addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logInstance.Error("Server failed to start", err)
		log.Fatalf("Server failed: %v", err)
	}
}

func initializeLogger() *logger.Logger {
	logInstance, err := logger.NewLogger("movie-service.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logInstance
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Environment Variables
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Printf("No .env file found or failed to load: %v", err)
		}
	}

	// Database connection
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatalf("DATABASE_URL not set in environment")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Try to find SQL file in multiple locations
	sqlPaths := []string{
		"database/import/database-dump.sql",    // From server root
		"./database/import/database-dump.sql",  // Relative to current dir
		"../database/import/database-dump.sql", // From import dir
		"database-dump.sql",                    // Fallback
	}

	var sqlFilePath string
	var sqlContent []byte
	var readErr error

	for _, path := range sqlPaths {
		if _, statErr := os.Stat(path); statErr == nil {
			sqlFilePath = path
			sqlContent, readErr = ioutil.ReadFile(path)
			if readErr == nil {
				log.Printf("Found SQL file at: %s", path)
				break
			}
		}
	}

	if sqlFilePath == "" || readErr != nil {
		log.Fatalf("Failed to find or read SQL file. Tried: %v", sqlPaths)
	}

	statements := strings.Split(string(sqlContent), ";\n")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		lines := strings.Split(stmt, "\n")
		var cleanedLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "--") {
				continue
			}
			cleanedLines = append(cleanedLines, trimmed)
		}

		cleanedStmt := strings.Join(cleanedLines, " ")
		if cleanedStmt == "" {
			continue
		}

		_, err := db.Exec(cleanedStmt)
		if err != nil {
			log.Printf("Failed to execute statement: %v\nStatement: %s\n", err, cleanedStmt)
			return
		}
		fmt.Printf("Executed: %s\n", cleanedStmt[:min(50, len(cleanedStmt))]+"...")
	}

	fmt.Println("SQL script execution completed.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

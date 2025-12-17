package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"	
	"log"		
	"os"
	"strings"	

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
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

	sqlFilePath := "database/import/database-dump.sql"
	if _, err := os.Stat(sqlFilePath); os.IsNotExist(err) {
		sqlFilePath = "database-dump.sql"
	}

	sqlContent, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("Failed to read SQL file %s: %v", sqlFilePath, err)
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

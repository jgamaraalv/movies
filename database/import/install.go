package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:password@127.0.0.1:5432/movies?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	sqlFilePath := "database/import/database-dump.sql"
	sqlContent, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatal("Failed to read SQL file:", err)
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

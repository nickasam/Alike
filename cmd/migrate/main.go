package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Database connection
	db, err := sql.Open("postgres", "host=localhost port=5432 user=alike_user dbname=alike_db sslmode=disable password=your_password")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	switch command {
	case "up":
		runMigrations(db, "up")
	case "down":
		runMigrations(db, "down")
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}

func runMigrations(db *sql.DB, direction string) {
	migrationsDir := "db/migrations"
	
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*_"+direction+".sql"))
	if err != nil {
		log.Fatalf("Failed to find migration files: %v", err)
	}

	for _, file := range files {
		fmt.Printf("Running migration: %s\n", file)
		
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file: %v", err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration: %v", err)
		}

		fmt.Printf("✓ Migration completed: %s\n", filepath.Base(file))
	}

	fmt.Println("All migrations completed successfully!")
}

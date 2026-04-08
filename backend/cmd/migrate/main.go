package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://mockstarket:mockstarket_dev@localhost:5432/mockstarket?sslmode=disable"
	}

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("warning: migrator close errors: src=%v db=%v", srcErr, dbErr)
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down|version]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("Migrations applied successfully.")

	case "down":
		steps := 1
		if len(os.Args) > 2 {
			if _, err := fmt.Sscanf(os.Args[2], "%d", &steps); err != nil {
				log.Fatalf("invalid steps %q: %v", os.Args[2], err)
			}
		}
		if err := m.Steps(-steps); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Printf("Rolled back %d migration(s).\n", steps)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("failed to get version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)

	default:
		fmt.Println("Usage: migrate [up|down|version]")
		os.Exit(1)
	}
}

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	db "be-dashboard-nba/internal/infrastructure/database"
	"be-dashboard-nba/internal/infrastructure/logger"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}
	logger.NewLogger()
	defer logger.Close()

	ctx := context.Background()
	database, err := db.NewDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if _, err = database.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS navi_schema_migrations (version TEXT PRIMARY KEY, applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		log.Fatal(err)
	}

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}
	migrationsDir := filepath.Join("database", "sqlite", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatal(err)
	}

	if direction == "down" {
		var version string
		err = database.QueryRowContext(ctx, `SELECT version FROM navi_schema_migrations ORDER BY version DESC LIMIT 1`).Scan(&version)
		if err != nil {
			log.Printf("No SQLite migration to rollback: %v", err)
			return
		}
		downSQL, readErr := os.ReadFile(filepath.Join(migrationsDir, version+".down.sql"))
		if readErr != nil {
			log.Fatal(readErr)
		}
		if _, err = database.ExecContext(ctx, string(downSQL)); err != nil {
			log.Fatal(err)
		}
		if _, err = database.ExecContext(ctx, `DELETE FROM navi_schema_migrations WHERE version=$1`, version); err != nil {
			log.Fatal(err)
		}
		log.Printf("SQLite migration rolled back: %s", version)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		version := strings.TrimSuffix(entry.Name(), ".up.sql")
		var applied int
		if err = database.QueryRowContext(ctx, `SELECT COUNT(*) FROM navi_schema_migrations WHERE version=$1`, version).Scan(&applied); err != nil {
			log.Fatal(err)
		}
		if applied > 0 {
			continue
		}
		migration, readErr := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
		if readErr != nil {
			log.Fatal(readErr)
		}
		if _, err = database.ExecContext(ctx, string(migration)); err != nil {
			log.Fatalf("SQLite migration %s failed: %v", version, err)
		}
		if _, err = database.ExecContext(ctx, `INSERT INTO navi_schema_migrations(version) VALUES($1)`, version); err != nil {
			log.Fatal(err)
		}
		log.Printf("SQLite migration applied: %s", version)
	}
	log.Printf("SQLite migrations applied to %s", os.Getenv("DB_SCHEMA"))
}

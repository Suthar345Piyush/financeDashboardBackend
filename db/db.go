// database related configuration / database connection

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

// function for connection to db

func Connect(dbPath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	// sqlite connection - wal(write ahead logging) mode for concurrent reads  and better concurrency support

	// using PRAGMA to modify operations in runtime

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_key=ON;",
		"PRAGMA busy_timeout=5000;",
	}

	for _, p := range pragmas {
		if _, err := database.Exec(p); err != nil {
			return nil, fmt.Errorf("pragma error (%s): %w", p, err)
		}
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Database Connected:", dbPath)

	return database, nil

}

// function to run migrations

func RunMigrations(database *sql.DB, migrationsDir string) error {

	entries, err := os.ReadDir(migrationsDir)

	if err != nil {
		return fmt.Errorf("Cannot read migrations dir: %w", err)
	}

	var files []string

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, filepath.Join(migrationsDir, e.Name()))
		}
	}

	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)

		if err != nil {
			return fmt.Errorf("cannot read migration %s: %w", f, err)
		}

		if _, err := database.Exec(string(content)); err != nil {
			return fmt.Errorf("migration failed %s: %w", f, err)
		}

		log.Println("Migration applied:", filepath.Base(f))

	}

	return nil

}

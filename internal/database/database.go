package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	ErrMigrationFailed = errors.New("an error occurred while processing database migration")
)

func getLastMigrationApplied(db *sql.DB) (string, error) {
	// Create the migrations table if it doesn't exist.
	// The `IF NOT EXISTS` clause is crucial for idempotency.
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		name TEXT PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return "", fmt.Errorf("%w: failed to create migrations table", ErrMigrationFailed)
	}

	// Query for the name of the last applied migration, ordered by applied_at timestamp.
	var lastMigration string
	row := db.QueryRow("SELECT name FROM migrations ORDER BY applied_at DESC LIMIT 1")
	err = row.Scan(&lastMigration)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No migrations have been applied yet.
			return "", nil
		}
		return "", fmt.Errorf("%w: failed to get last applied migration", ErrMigrationFailed)
	}

	return lastMigration, nil
}

// Gets a list of migrations that have been created after the last migration
func getUnappliedMigrations(lastApplied string) ([]string, error) {
	// For this to work, you need a 'db/migrations' directory
	// and your files must be named like 'YYYYMMDDhhmmss_name.sql'
	dir := "dbmigration"
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: could not read migrations directory: %v", ErrMigrationFailed, err)
	}

	var unapplied []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// The migration file name itself is the timestamp.
		fileName := strings.TrimSuffix(file.Name(), ".sql")

		if lastApplied == "" || fileName > lastApplied {
			unapplied = append(unapplied, filepath.Join(dir, file.Name()))
		}
	}

	sort.Strings(unapplied) // Ensure migrations are applied in sequential order.
	return unapplied, nil
}

func initializeDatabase(db *sql.DB) error {
	lastMigration, err := getLastMigrationApplied(db)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}

	migrations, err := getUnappliedMigrations(lastMigration)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationFailed, err)
	}

	if len(migrations) == 0 {
		fmt.Println("No new migrations to apply.")
		return nil
	}

	for _, migrationFile := range migrations {
		fmt.Printf("Applying migration: %s\n", migrationFile)

		// Read the SQL content from the file
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("%w: failed to read migration file %s: %v", ErrMigrationFailed, migrationFile, err)
		}

		// Use a transaction to ensure all commands in a migration file either succeed or fail as a unit.
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("%w: failed to begin transaction for %s: %v", ErrMigrationFailed, migrationFile, err)
		}
		defer tx.Rollback()

		if _, err := tx.Exec(string(content)); err != nil {
			if liteErr, ok := err.(*sqlite.Error); ok {
				code := liteErr.Code()
				if code == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
					// This migration has already been applied, skip it.
					continue
				}
			}
			return fmt.Errorf("%w: failed to execute migration %s: %v", ErrMigrationFailed, migrationFile, err)
		}

		// Get the migration file name (e.g., '20250816123000_create_users_table')
		name := strings.TrimSuffix(filepath.Base(migrationFile), ".sql")

		// Record the applied migration in the migrations table
		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", name); err != nil {
			if liteErr, ok := err.(*sqlite.Error); ok {
				code := liteErr.Code()
				if code == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
					// This migration has already been applied, skip it.
					if err := tx.Commit(); err != nil {
						return fmt.Errorf("%w: failed to commit transaction for %s: %v", ErrMigrationFailed, migrationFile, err)
					}
					continue
				}
			}
			return fmt.Errorf("%w: failed to record migration %s: %v", ErrMigrationFailed, migrationFile, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("%w: failed to commit transaction for %s: %v", ErrMigrationFailed, migrationFile, err)
		}
	}

	return nil
}

func CreateDatabaseConnection(filePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s", filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", filePath, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database at %s: %w", filePath, err)
	}

	if err := initializeDatabase(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: could not initialize database", err)
	}

	return db, nil
}

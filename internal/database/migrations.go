package database

import (
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(db *gorm.DB) error {
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`).Error; err != nil {
		return err
	}

	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(file, ".sql")

		var alreadyApplied bool
		if err := db.Raw(
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = ?)`,
			version,
		).Scan(&alreadyApplied).Error; err != nil {
			return err
		}

		if alreadyApplied {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + file)
		if err != nil {
			return err
		}

		log.Printf("rodando migration: %s", file)

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(content)).Error; err != nil {
				return fmt.Errorf("erro na migration %s: %w", file, err)
			}

			if err := tx.Exec(
				`INSERT INTO schema_migrations (version) VALUES (?)`,
				version,
			).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// internal/db/migrate.go
package db

import (
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"path/filepath"
)

func RunMigrations(db *sqlx.DB) {
	migrationsPath := "./migrations"
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		log.Fatal("Failed to read migrations: ", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal("Failed to read migration file ", file, ": ", err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal("Failed to run migration ", file, ": ", err)
		}

		log.Println("✅ Ran migration:", filepath.Base(file))
	}
}
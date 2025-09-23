package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"brutal/internal/config"
)

var DB *sqlx.DB

func InitDB(cfg *config.Config) {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	// Open connection
	var err error
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to connect to DB: ", err)
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ DB ping failed: ", err)
	}
	log.Println("✅ Connected to PostgreSQL!")

}
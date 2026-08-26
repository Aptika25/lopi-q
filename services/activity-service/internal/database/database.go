package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"activity-service/internal/config"
)

func NewDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("⚠️ PostgreSQL Connection Error (activity-service): %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Printf("⚠️ PostgreSQL Ping Failed (activity-service): %v", err)
		return nil
	}

	log.Printf("✅ Connected to PostgreSQL database: %s", cfg.DBName)
	return db
}

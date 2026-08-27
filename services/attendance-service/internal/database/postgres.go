package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"attendance-service/internal/config"
)

func NewDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to connect to postgres: %v. Using fallback in-memory.", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Printf("⚠️ Warning: Failed to ping postgres: %v. Using fallback in-memory.", err)
		return nil
	}

	return db
}

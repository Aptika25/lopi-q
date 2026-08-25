package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"auth-service/internal/config"

	_ "github.com/lib/pq"
)

func NewDB(cfg *config.Config) *sql.DB {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Printf("[Auth-DB Warning] Unable to open PostgreSQL connection (%v). Dynamic fallback active.", err)
		return nil
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("[Auth-DB Warning] Could not ping PostgreSQL at %s:%s (%v).", cfg.DBHost, cfg.DBPort, err)
		return nil
	}
	log.Printf("Successfully connected to PostgreSQL database [%s]", cfg.DBName)

	return db
}

func NewUserDB(cfg *config.Config) *sql.DB {
	conn := fmt.Sprintf(
		"host=%s port=%s user=user_lopiq_user password=lopiquserPassword@2k26# dbname=db_lopiq_user sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil
	}
	log.Printf("Successfully connected to PostgreSQL database [db_lopiq_user] for 2FA sync")
	return db
}

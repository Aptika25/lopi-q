package config

import "os"

type Config struct {
	Port       string
	DBType     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBPath     string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres_apps"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "user_lopiq_user"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "lopiquserPassword@2k26#"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "db_lopiq_user"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		candidates := []string{
			"/app/data/users.json",
			"data/users.json",
			"../data/users.json",
			"../../data/users.json",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				dbPath = c
				break
			}
		}
		if dbPath == "" {
			dbPath = "/app/data/users.json"
		}
	}
	return &Config{
		Port:       port,
		DBType:     "postgres",
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPass,
		DBName:     dbName,
		DBSSLMode:  "disable",
		DBPath:     dbPath,
	}
}

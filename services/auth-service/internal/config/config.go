package config

import "os"

type Config struct {
	Port       string
	JWTSecret  string
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
		port = "50051"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "lopiq_super_secret_jwt_key_2026"
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
		dbUser = "user_lopiq_auth"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "lopiqauthPassword@2k26#"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "db_lopiq_auth"
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
		JWTSecret:  jwtSecret,
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

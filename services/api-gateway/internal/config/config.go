package config

import "os"

type Config struct {
	Port               string
	JWTSecret          string
	AuthServiceURL     string
	UserServiceURL     string
	PresensiServiceURL string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "lopiq_super_secret_jwt_key_2026"
	}
	return &Config{
		Port:               port,
		JWTSecret:          jwtSecret,
		AuthServiceURL:     "localhost:50051",
		UserServiceURL:     "localhost:50052",
		PresensiServiceURL: "localhost:50053",
	}
}

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
	// Tokens are signed by pkg/jwt, which reads JWT_SECRET itself; this copy is informational only.
	jwtSecret := os.Getenv("JWT_SECRET")
	return &Config{
		Port:               port,
		JWTSecret:          jwtSecret,
		AuthServiceURL:     "localhost:50051",
		UserServiceURL:     "localhost:50052",
		PresensiServiceURL: "localhost:50053",
	}
}

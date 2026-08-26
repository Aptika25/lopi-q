package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port       string
	DBType     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() *Config {
	loadDotEnv()

	return &Config{
		Port:       getEnv("PORT", "50053"),
		DBType:     getEnv("DB_TYPE", "postgres"),
		DBHost:     getEnv("DB_HOST", "postgres_apps"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user_lopiq_activity"),
		DBPassword: getEnv("DB_PASSWORD", "lopiqactivityPassword@2k26#"),
		DBName:     getEnv("DB_NAME", "db_lopiq_activity"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func loadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

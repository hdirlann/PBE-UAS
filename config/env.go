package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Env holds environment configuration
type Env struct {
	AppPort string

	JWTSecret       string
	JWTExpiresHours int

	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDatabase string
	PGSSLMode  string

	MongoURI string
	MongoDB  string
}

var (
	cfg  *Env
	once sync.Once
)

// LoadEnv loads .env once and returns the env singleton
func LoadEnv() *Env {
	once.Do(func() {
		_ = godotenv.Load() // ignore error; allow OS env override

		cfg = &Env{
			AppPort:         getEnv("APP_PORT", "8080"),
			JWTSecret:       getEnv("JWT_SECRET", "change-this-secret"),
			JWTExpiresHours: getEnvInt("JWT_EXPIRES_HOURS", 24),

			PGHost:     getEnv("PG_HOST", "localhost"),
			PGPort:     getEnv("PG_PORT", "5432"),
			PGUser:     getEnv("PG_USER", "postgres"),
			PGPassword: getEnv("PG_PASSWORD", ""),
			PGDatabase: getEnv("PG_DATABASE", "prestasi_db"),
			PGSSLMode:  getEnv("PG_SSLMODE", "disable"),

			MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			MongoDB:  getEnv("MONGO_DBNAME", "prestasi_db"),
		}

		if cfg.JWTSecret == "" {
			log.Println("WARNING: JWT_SECRET is empty; change for production")
		}
	})
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("invalid int for %s, using default %d", key, fallback)
	}
	return fallback
}

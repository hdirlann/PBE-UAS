package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppPort    string
	JwtSecret  string
	ApiKey     string

	// Postgres
	PgHost     string
	PgPort     string
	PgUser     string
	PgPassword string
	PgDatabase string

	// Mongo
	MongoURI  string
	MongoDB   string
}

var Cfg Env

// LoadEnv reads .env into Cfg
func LoadEnv() {
	_ = godotenv.Load()

	Cfg = Env{
		AppPort:    getEnv("APP_PORT", "8080"),
		JwtSecret:  getEnv("JWT_SECRET", "secret"),
		ApiKey:     getEnv("API_KEY", ""),

		PgHost:     getEnv("PG_HOST", "localhost"),
		PgPort:     getEnv("PG_PORT", "5432"),
		PgUser:     getEnv("PG_USER", "postgres"),
		PgPassword: getEnv("PG_PASSWORD", ""),
		PgDatabase: getEnv("PG_DATABASE", "appdb"),

		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGO_DBNAME", "appdb"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Init() {
	LoadEnv()
	if Cfg.JwtSecret == "" {
		log.Println("WARNING: JWT_SECRET empty")
	}
}

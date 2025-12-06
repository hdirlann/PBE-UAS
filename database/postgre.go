package database

import (
	"database/sql"
	"fmt"
	"time"

	"clean-arch/config" // sesuaikan jika module name berbeda

	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

// ConnectPostgres connects and sets the package-global PostgresDB.
// Returns error on failure.
//
// Usage (classic): call ConnectPostgres(env) and repositories can use database.PostgresDB.
func ConnectPostgres(env *config.Env) error {
	if env == nil {
		return fmt.Errorf("config env is nil")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env.PGHost, env.PGPort, env.PGUser, env.PGPassword, env.PGDatabase, env.PGSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// optional tuning
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return err
	}

	// assign to global variable used across repo
	PostgresDB = db
	return nil
}

// ConnectPostgresReturn is an alternative that returns the *sql.DB in case you prefer DI.
// It also assigns PostgresDB global for backwards compatibility.
func ConnectPostgresReturn(env *config.Env) (*sql.DB, error) {
	if err := ConnectPostgres(env); err != nil {
		return nil, err
	}
	return PostgresDB, nil
}

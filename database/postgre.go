package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var PostgresDB *sql.DB

// ConnectPostgres connects using env: PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DATABASE
func ConnectPostgres() error {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DATABASE")

	if host == "" { host = "localhost" }
	if port == "" { port = "5432" }
	if user == "" { user = "postgres" }
	if dbname == "" { dbname = "appdb" }

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// simple ping
	if err := db.Ping(); err != nil {
		return err
	}

	PostgresDB = db
	log.Printf("connected to postgres: %s:%s/%s\n", host, port, dbname)
	return nil
}

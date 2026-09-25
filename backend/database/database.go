package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "ordet_bank"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return runMigrations(db)
}

func runMigrations(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		balance NUMERIC(14,2) NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		account_id INTEGER NOT NULL REFERENCES accounts(id),
		type TEXT NOT NULL,
		amount NUMERIC(14,2) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS transfers (
		id SERIAL PRIMARY KEY,
		from_account INTEGER NOT NULL REFERENCES accounts(id),
		to_account INTEGER NOT NULL REFERENCES accounts(id),
		amount NUMERIC(14,2) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now()
	);
	`
	_, err := db.Exec(schema)
	return err
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
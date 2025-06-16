package adapters

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func NewPostgresConnection() (*sqlx.DB, error) {
	POSTGRES_USER := getEnvOrDefault("POSTGRES_USER", "user")
	POSTGRES_PASSWORD := getEnvOrDefault("POSTGRES_PASSWORD", "secret")
	POSTGRES_ADDR := getEnvOrDefault("POSTGRES_ADDR", "localhost")
	POSTGRES_DB := getEnvOrDefault("POSTGRES_DB", "wild-workout")

	dsn := "postgres://" +
		POSTGRES_USER + ":" +
		POSTGRES_PASSWORD + "@" +
		POSTGRES_ADDR + "/" +
		POSTGRES_DB + "?sslmode=disable"

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to PostgreSQL")
	}

	return db, nil
}

func MustNewPostgresConnection() *sqlx.DB {
	db, err := NewPostgresConnection()
	if err != nil {
		panic(err)
	}

	return db
}

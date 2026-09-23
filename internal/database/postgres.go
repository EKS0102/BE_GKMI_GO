package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Connect() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../../.env")
		if err != nil {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbName,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return conn, nil
}

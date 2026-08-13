package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func DbConnect() (*pgxpool.Pool, context.Context, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading .env file: %w", err)
	}

	connString := fmt.Sprintf(
		"postgres://postgres:%s@%s:%s/%s",
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, ctx, nil
}

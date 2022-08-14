package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// NewPostgres is a main database
func NewPostgres() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		panic(err)
	}

	return dbpool
}

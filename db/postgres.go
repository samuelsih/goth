package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// NewPostgres is a main database
func NewPostgres() *pgxpool.Pool {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s", 
	os.Getenv("YSQL_HOST"), 
	os.Getenv("YSQL_PORT"), 
	os.Getenv("YSQL_USER"),
	os.Getenv("YSQL_PASSWORD"),
	os.Getenv("YSQL_DB"),
	os.Getenv("YSQL_SSL"),
	os.Getenv("YSQL_CERT"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		panic(err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		panic(err)
	}

	log.Println("Postgres ready!")

	return dbpool
}

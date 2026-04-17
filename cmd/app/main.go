package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()
	connString := "postgres://postgres:123456@localhost:5432/postgres?sslmode=disable"
	pool, err := connectDB(ctx, connString)
	if err != nil {
		errors.Wrap(entities.ErrInternal, "failed to connect to the database")
	}
	defer pool.Close()

	// добавить возвращаемое значение
	//postgres.NewStorage(pool)

}

func connectDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, connString)
}

func addEnv() string {
	//connString := "postgres://{DB_USER}:{DB_PASSWORD}@localhost:5432/{DB_NAME}?sslmode=disable"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPassword, dbName)
	return connString
}

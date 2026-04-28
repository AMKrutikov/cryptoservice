package main

import (
	"fmt"
	"os"

	coingecko "github.com/AMKrutikov/cryptoservice/internal/adapter/provider"
	"github.com/AMKrutikov/cryptoservice/internal/adapter/storage/postgres"
	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/AMKrutikov/cryptoservice/internal/port/http/public"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	// нужно настроить ожидание сигнала Ctrl+C

	// storage
	connString, err := connPostgres()
	if err != nil {
		fmt.Println(err)
	}

	storage, err := postgres.NewStorage(connString)
	if err != nil {
		fmt.Printf("Storage error: %v\n", err)
		return
	}
	defer storage.Close()
	//

	//provider
	provider, err := coingecko.NewProviderClient("x-cg-demo-api-key", "CG-SdGMn7C5Rv2F4hTMwLdJ1Pk6")
	if err != nil {
		fmt.Printf("Provider error: %v\n", err)
		return
	}
	//

	//service
	service, err := cases.NewService(provider, storage)
	if err != nil {
		fmt.Printf("Service error: %v\n", err)
		return
	}

	serverRun := public.NewServer(service)                 ///
	if err := serverRun.StartServer(":9091"); err != nil { ///
		fmt.Println("Error start server") ///
	}

}

func connPostgres() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", errors.Wrapf(entities.ErrInvalidParam, "error loading .env file values:%v", err)
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	//connString := "postgres://postgres:123456@localhost:5432/postgres?sslmode=disable"
	//connString := "postgres://{DB_USER}:{DB_PASSWORD}@localhost:5432/{DB_NAME}?sslmode=disable"
	connString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPassword, dbName)
	return connString, nil
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	coingecko "github.com/AMKrutikov/cryptoservice/internal/adapter/provider"
	"github.com/AMKrutikov/cryptoservice/internal/adapter/storage/postgres"
	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/entities"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	// нужно настроить ожидание сигнала Ctrl+C

	ctx := context.Background()
	titles := []string{"solana", "tron", "bitcoin", "rain", "ripple", "zcash", "dai"}
	//titles := []string{"xrp"}

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

	provider, err := coingecko.NewProviderClient("x-cg-demo-api-key", "CG-SdGMn7C5Rv2F4hTMwLdJ1Pk6")
	if err != nil {
		fmt.Printf("Provider error: %v\n", err)
		return
	}

	service, err := cases.NewService(provider, storage)
	if err != nil {
		fmt.Printf("Service error: %v\n", err)
		return
	}

	result, err := service.GetLastRates(ctx, titles)
	if err != nil {
		fmt.Println(err)
	}

	// err = service.ActualizeRates(ctx)
	// if err != nil {
	// 	fmt.Println("error ActualizeRates", err)
	// }

	// result, err := service.GetAgregetedRates(ctx, titles, "a")
	// if err != nil {
	// 	fmt.Println("error AcGetAgregetedRates", err)
	// }

	for _, elem := range result {
		fmt.Printf("Монета: %s Цена: %.4f Время:%s\n", elem.Title(), elem.Price(), elem.ActualAt().Format(time.Stamp))
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

package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AMKrutikov/cryptoservice/internal/adapter/config"
	coingecko "github.com/AMKrutikov/cryptoservice/internal/adapter/provider"
	"github.com/AMKrutikov/cryptoservice/internal/adapter/storage/postgres"
	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/AMKrutikov/cryptoservice/internal/port/http/public"
	"github.com/robfig/cron/v3"
)

type Application struct {
	cnfg *config.Config

	provider cases.CryptoProvider
	storage  cases.CryptoStorage
	service  port.Service
	cron     *cron.Cron

	publicPort *public.Server
}

func NewApplication(cnfg *config.Config) *Application {
	return &Application{
		cnfg: cnfg,
	}
}

func (app *Application) Run(ctx context.Context) {
	app.initCryptoProvider()
	app.initStorage()
	app.initService()
	app.initPublicHTTPPort()
	app.initCron(ctx)

	go func() {
		if err := app.startHTTPPublic(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v", err)
		}
	}()

	fmt.Println("Application is running...")

}

func (app *Application) Stop() {
	fmt.Println("Shutting down application...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if app.cron != nil {
		fmt.Println("Stopping cron scheduler...")
		cronCtxStop := app.cron.Stop()
		select {
		case <-cronCtxStop.Done():
			fmt.Println("Cron scheduler shutdown completed")
		case <-shutdownCtx.Done():
			fmt.Println("Cron scheduler stopped timeout")
		}
	}

	if app.publicPort != nil {
		fmt.Println("Stopping httpServer...")
		if err := app.publicPort.Stop(shutdownCtx); err != nil {
			fmt.Printf("HTTP server shutdown error: %v\n", err)
		} else {
			fmt.Println("HTTP server shutdown completed")
		}
	}

	if app.storage != nil {
		fmt.Println("Closing database connections...")
		app.storage.Close()
		fmt.Println("Database shutdown completed")
	}

}

func (app *Application) initCryptoProvider() {
	apiKey := app.cnfg.CoingeckoAPIKey()
	client, err := coingecko.NewProviderClient(apiKey)
	if err != nil {
		panic(err)
	}

	app.provider = client
}

func (app *Application) initStorage() {
	storageType := app.cnfg.StorageType()
	connString := app.cnfg.StorageConnectionString(storageType)

	storage, err := postgres.NewStorage(connString)
	if err != nil {
		panic(err)
	}

	app.storage = storage
}

func (app *Application) initService() {
	service, err := cases.NewService(app.provider, app.storage)
	if err != nil {
		panic(err)
	}
	app.service = service
}

func (app *Application) initPublicHTTPPort() {
	server := public.NewServer(app.service)
	app.publicPort = server
}

func (app *Application) startHTTPPublic() error {
	port := app.cnfg.PublicHTTPPort()
	address := ":" + port
	return app.publicPort.StartServer(address)

}

func (app *Application) initCron(ctx context.Context) {

	app.cron = cron.New(cron.WithSeconds())
	Schedule := "0 */1 * * * *"

	_, err := app.cron.AddFunc(Schedule, func() {
		ctxTime, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		if err := app.service.ActualizeRates(ctxTime); err != nil {
			fmt.Printf("Cron execution error: %v\n", err)
			return
		}
		fmt.Println("Cron: Crypto rates successfully actualized.")
	})
	if err != nil {
		panic(fmt.Errorf("failed to start the cron job: %w\n", err))
	}

	app.cron.Start()
}

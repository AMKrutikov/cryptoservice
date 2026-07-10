package application

import (
	"context"
	"fmt"
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

	if err := app.startHTTPPublic(); err != nil {
		panic(err)
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
	Shedule := "0 */1 * * * *"

	_, err := app.cron.AddFunc(Shedule, func() {
		ctxTime, cancel := context.WithTimeout(ctx, time.Second*30)
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

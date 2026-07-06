package application

import (
	"github.com/AMKrutikov/cryptoservice/internal/adapter/config"
	coingecko "github.com/AMKrutikov/cryptoservice/internal/adapter/provider"
	"github.com/AMKrutikov/cryptoservice/internal/adapter/storage/postgres"
	"github.com/AMKrutikov/cryptoservice/internal/cases"
	"github.com/AMKrutikov/cryptoservice/internal/port"
	"github.com/AMKrutikov/cryptoservice/internal/port/http/public"
)

type Application struct {
	cnfg *config.Config

	provider cases.CryptoProvider
	storage  cases.CryptoStorage
	service  port.Service

	publicPort *public.Server
}

func NewApplication(cnfg *config.Config) *Application {
	return &Application{
		cnfg: cnfg,
	}
}

func (app *Application) Run() {
	app.initCryptoProvider()
	app.initStorage()
	app.initService()
	app.initPublicHTTPPort()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// app.initActualizeInterval(ctx)

	if err := app.startHTTPPublic(); err != nil {
		panic(err)
	}
}

func (app *Application) initCryptoProvider() {
	authKey := app.cnfg.CoingeckoAuthKey()
	apiKey := app.cnfg.CoingeckoAPIKey()
	client, err := coingecko.NewProviderClient(authKey, apiKey)
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

// func (app *Application) initActualizeInterval(ctx context.Context) {
// 	interval := app.cnfg.ActualizeInterval()

// 	ticker := time.NewTicker(interval)

// 	go func() {
// 		defer ticker.Stop()
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case <-ticker.C:
// 				app.service.ActualizeRates(ctx)
// 			}
// 		}
// 	}()

// }

// cronJob

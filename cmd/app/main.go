package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/AMKrutikov/cryptoservice/internal/adapter/config"
	"github.com/AMKrutikov/cryptoservice/pkg/application"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	path := os.Getenv("PATH_CONFIG")
	if strings.TrimSpace(path) == "" {
		path = "config/cryptoservice.yaml"
	}

	cnfg := config.NewConfig(path)
	app := application.NewApplication(cnfg)
	app.Run(ctx)

	<-ctx.Done()
	fmt.Println("Main: Received shutdown signal")

	app.Stop()
	fmt.Println("Main: Application stopped cleanly")

}

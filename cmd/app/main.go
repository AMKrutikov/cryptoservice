package main

import (
	"os"
	"strings"

	"github.com/AMKrutikov/cryptoservice/internal/adapter/config"
	"github.com/AMKrutikov/cryptoservice/pkg/application"
)

func main() {
	path := os.Getenv("PATH_CONFIG")
	if strings.TrimSpace(path) == "" {
		path = "config/cryptoservice.yaml"
	}

	cnfg := config.NewConfig(path)
	app := application.NewApplication(cnfg)
	app.Run()

}

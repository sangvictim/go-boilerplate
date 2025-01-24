package main

import (
	"fmt"
	"go-api-fiber/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Validate: validator,
		Log:      log,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt32("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

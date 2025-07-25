package main

import (
	"github.com/thoriqwildan/aino-medical-be/db/seed"
	"github.com/thoriqwildan/aino-medical-be/internal/config"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/middleware"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := config.NewValidator(viperConfig)
	jwtMiddleware := middleware.NewMiddlewareConfig(viperConfig, app)

	config.Bootstrap(&config.BootstrapConfig{
		DB: db,
		App: app,
		Log: log,
		Validate: validator,
		Config: viperConfig,
		JWT: jwtMiddleware,
	})

	seeding := viperConfig.GetBool("SEED")
	if seeding {
		seed.RunAllSeeders(db)
	}

	webPort := ":" + viperConfig.GetString("WEB_PORT")
	err := app.Listen(webPort)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	log.Info("Server is running on port ", webPort)
}
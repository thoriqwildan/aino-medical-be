package main

import (
	"github.com/thoriqwildan/aino-medical-be/db/seed"
	"github.com/thoriqwildan/aino-medical-be/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	seed.RunAllSeeders(db)
}

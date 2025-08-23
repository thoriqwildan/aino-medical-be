package main

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/thoriqwildan/aino-medical-be/internal/config"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, logger)
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("load tz: %v", err)
	}
	c := cron.New(
		cron.WithLocation(location),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)
	logger.Info("Starting cron job")

	c.AddFunc("0 0 * * *", func() {
		logger.Info("Do Cron Job Check Employee Pro Rate")
		var employees []*entity.Employee
		if err := db.Find(&employees).Error; err != nil {
			logger.Error("Error when find the all employee", err)
			log.Fatal(err)
		}

		for _, employee := range employees {
			now := time.Now()

			join := employee.JoinDate
			if join.Location() != now.Location() {
				now = now.In(join.Location())
			}

			oneYearElapsed := !now.Before(join.AddDate(1, 0, 0))

			if !join.IsZero() && oneYearElapsed {
				if err := db.Model(&entity.Employee{}).
					Where("id = ?", employee.ID).
					Update("pro_rate", 0).Error; err != nil {
					logger.Errorf("set pro_rate=0 for employee %d: %v", employee.ID, err)
				}
			}

		}
	}) // every 00:00

	c.Run()
}

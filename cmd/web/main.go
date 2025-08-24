package main

import (
	"context"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thoriqwildan/aino-medical-be/internal/config"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/middleware"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"

	_ "github.com/thoriqwildan/aino-medical-be/docs"
)

// @title Aino Medical API
// @version 1.0
// @description This is a sample swagger for Fiber
// @host localhost:3000
// @BasePath /

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @securityschemes.bearer.name Authorization
// @securityschemes.bearer.in header
// @securityschemes.bearer.type http
// @securityschemes.bearer.scheme bearer
// @securityschemes.bearer.description Type "Bearer" followed by a space and JWT token.
func main() {
	// ===== Context for graceful shutdown =====
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ===== Bootstrap Component =====
	viperConfig := config.NewViper()
	logg := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, logg)
	validator := config.NewValidator(viperConfig)
	jwtMiddleware := middleware.NewMiddlewareConfig(viperConfig, app)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      logg,
		Validate: validator,
		Config:   viperConfig,
		JWT:      jwtMiddleware,
	})

	app.Static("/swagger", "./docs")

	if runtime.GOOS == "windows" {
		app.Get("/docs", func(ctx *fiber.Ctx) error {
			html, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL:  ctx.BaseURL() + "/swagger/swagger.json",
				DarkMode: true,
				Theme:    scalar.ThemeKepler,
				Layout:   scalar.LayoutModern,
			})
			if err != nil {
				logg.Error("Failed to generate API reference HTML: " + err.Error())
				return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to generate API reference HTML")
			}
			return ctx.Type("html").SendString(html)
		})
	} else {
		app.Get("/docs", func(ctx *fiber.Ctx) error {
			html, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL:  "./docs/swagger.json",
				DarkMode: true,
				Theme:    scalar.ThemeKepler,
				Layout:   scalar.LayoutModern,
			})
			if err != nil {
				log.Error("Failed to generate API reference HTML: " + err.Error())
				return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to generate API reference HTML")
			}
			return ctx.Type("html").SendString(html)
		})
	}

	// Health endpoints
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("OK") })
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("healthy") })

	// ====== Initialize & start Cron ======
	cjobs, errCronJobs := config.NewCronJobs("Asia/Jakarta")
	if errCronJobs != nil {
		logg.Error(errCronJobs)
	}
	cjobs.Add("reset-benefit-daily", config.NewCronJob(
		"reset-benefit-daily",
		func() error {
			start := time.Now()
			logg.Info("[cron] start reset-benefit-daily")

			err := helper.ResetBenefitProRateDaily(db)
			if err != nil {
				logg.Error(err.Error())
				return err
			}
			logg.Infof("[cron] done in %s", time.Since(start))
			return nil
		},
		"0 0 * * *",
	))
	cjobs.Add("reset-patient-benefit-daily", config.NewCronJob(
		"reset-patient-benefit-daily",
		func() error {
			start := time.Now()
			logg.Info("[cron] reset-patient-benefit-daily")

			err := helper.ResetPatientBenefitRemainingPlafondDaily(db)
			if err != nil {
				logg.Error(err.Error())
				return err
			}
			logg.Infof("[cron] done in %s", time.Since(start))
			return nil
		},
		"0 0 * * *",
	))
	cjobs.Run()

	serverErr := make(chan error, 1)
	go func() {
		webPort := ":" + viperConfig.GetString("WEB_PORT")
		logg.Infof("Server starting on %s", webPort)
		serverErr <- app.Listen(webPort)
	}()

	select {
	case <-ctx.Done():
		logg.Warn("Shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			logg.Error("Server error: ", err)
		}
	}

	shutdownCronDone := make(chan struct{})
	go func() {
		cjobs.Stop()
		close(shutdownCronDone)
	}()
	select {
	case <-shutdownCronDone:
	case <-time.After(2 * time.Minute):
		logg.Warn("Cron shutdown timeout; continue…")
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shCtx); err != nil {
		logg.Error("Fiber shutdown error: ", err)
	}

	logg.Info("Shutdown complete")
}

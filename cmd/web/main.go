package main

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thoriqwildan/aino-medical-be/db/seed"
	"github.com/thoriqwildan/aino-medical-be/internal/config"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/middleware"

	_ "github.com/thoriqwildan/aino-medical-be/docs"
)

// @title Aino Medical API
// @version 1.0
// @description This is a sample swagger for Fiber
// @host 192.168.74.32:3000
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

// @securityschemes.apiKey ApiKeyAuth
// @in header
// @name X-API-Key
// @description Enter your API Key as X-API-Key header value
func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := config.NewValidator(viperConfig)
	jwtMiddleware := middleware.NewMiddlewareConfig(viperConfig, app)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

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

	app.Get("/docs", func(ctx *fiber.Ctx) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			DarkMode: true,
			Theme: scalar.ThemeKepler,
			Layout: scalar.LayoutModern,
		})
		if err != nil {
			log.Error("Failed to generate API reference HTML: " + err.Error())
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to generate API reference HTML")
		}
		return ctx.Type("html").SendString(html)
	})

	webPort := ":" + viperConfig.GetString("WEB_PORT")
	err := app.Listen(webPort)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	log.Info("Server is running on port ", webPort)
}
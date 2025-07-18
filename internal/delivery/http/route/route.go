package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/http"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/middleware"
)

type RouteConfig struct {
	App *fiber.App
	JWT *middleware.MiddlewareConfig
	UserController *http.UserController
}

func (rc *RouteConfig) Setup() {
	rc.GeneralRoutes()
	rc.ProtectedRoutes()
}

func (rc *RouteConfig) GeneralRoutes() {
	general := rc.App.Group("/api/v1")
	general.Post("/auth/register", rc.UserController.Register)
	general.Post("/auth/login", rc.UserController.Login)
}

func (rc *RouteConfig) ProtectedRoutes() {
	protected := rc.App.Group("/api/v1/coba", rc.JWT.JWTProtected())
	protected.Get("/test", rc.UserController.GetTest)
	// Add more protected routes here
}

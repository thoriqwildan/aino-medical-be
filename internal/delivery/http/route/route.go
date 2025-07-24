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
	TransactionTypeController *http.TransactionTypeController
	PlanTypeController *http.PlanTypeController
}

func (rc *RouteConfig) Setup() {
	rc.GeneralRoutes()
	rc.ProtectedRoutes()
	rc.TransactionTypeRotes()
	rc.PlanTypeRoutes()
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

func (rc *RouteConfig) TransactionTypeRotes() {
	transactionType := rc.App.Group("/api/v1/transaction-types", rc.JWT.JWTProtected())
	transactionType.Post("/", rc.TransactionTypeController.Create)
	transactionType.Get("/:id", rc.TransactionTypeController.GetById)
	transactionType.Get("/", rc.TransactionTypeController.Get)
	transactionType.Put("/:id", rc.TransactionTypeController.Update)
	transactionType.Delete("/:id", rc.TransactionTypeController.Delete)
}

func (rc *RouteConfig) PlanTypeRoutes() {
	planType := rc.App.Group("/api/v1/plan-types", rc.JWT.JWTProtected())
	planType.Post("/", rc.PlanTypeController.Create)
	planType.Get("/:id", rc.PlanTypeController.GetById)
	planType.Get("/", rc.PlanTypeController.Get)
}

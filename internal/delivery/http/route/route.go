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
	LimitationTypeController *http.LimitationTypeController
	BenefitController *http.BenefitController
	DepartmentController *http.DepartmentController
	EmployeeController *http.EmployeeController
}

func (rc *RouteConfig) Setup() {
	rc.GeneralRoutes()
	rc.ProtectedRoutes()
	rc.TransactionTypeRotes()
	rc.PlanTypeRoutes()
	rc.LimitationTypeRoutes()
	rc.BenefitRoutes()
	rc.DepartmentRoutes()
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
	planType.Put("/:id", rc.PlanTypeController.Update)
	planType.Delete("/:id", rc.PlanTypeController.Delete)
}

func (rc *RouteConfig) LimitationTypeRoutes() {
	limitationType := rc.App.Group("/api/v1/limitation-types", rc.JWT.JWTProtected())
	limitationType.Post("/", rc.LimitationTypeController.Create)
	limitationType.Get("/:id", rc.LimitationTypeController.GetById)
	limitationType.Get("/", rc.LimitationTypeController.GetAll)
	limitationType.Put("/:id", rc.LimitationTypeController.Update)
	limitationType.Delete("/:id", rc.LimitationTypeController.Delete)
}

func (rc *RouteConfig) BenefitRoutes() {
	benefit := rc.App.Group("/api/v1/benefits", rc.JWT.JWTProtected())
	benefit.Post("/", rc.BenefitController.Create)
	benefit.Get("/:id", rc.BenefitController.GetById)
	benefit.Get("/", rc.BenefitController.GetAll)
	benefit.Put("/:id", rc.BenefitController.Update)
	benefit.Delete("/:id", rc.BenefitController.Delete)
}

func (rc *RouteConfig) DepartmentRoutes() {
	department := rc.App.Group("/api/v1/departments", rc.JWT.JWTProtected())
	department.Post("/", rc.DepartmentController.Create)
	department.Get("/:id", rc.DepartmentController.GetById)
	department.Get("/", rc.DepartmentController.GetAll)
	department.Put("/:id", rc.DepartmentController.Update)
	department.Delete("/:id", rc.DepartmentController.Delete)
}

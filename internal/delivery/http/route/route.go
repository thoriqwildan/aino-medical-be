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
	FamilyMemberController *http.FamilyMemberController
	ClaimController *http.ClaimController
}

func (rc *RouteConfig) Setup() {
	rc.GeneralRoutes()
	rc.ProtectedRoutes()
	rc.TransactionTypeRotes()
	rc.PlanTypeRoutes()
	rc.LimitationTypeRoutes()
	rc.BenefitRoutes()
	rc.DepartmentRoutes()
	rc.EmployeeRoutes()
	rc.FamilyMemberRoutes()
	rc.ClaimRoutes()
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

func (rc *RouteConfig) EmployeeRoutes() {
	employee := rc.App.Group("/api/v1/employees", rc.JWT.JWTProtected())
	employee.Post("/", rc.EmployeeController.Create)
	employee.Get("/:id", rc.EmployeeController.GetByID)
	employee.Get("/", rc.EmployeeController.GetAll)
	employee.Put("/:id", rc.EmployeeController.Update)
	employee.Delete("/:id", rc.EmployeeController.Delete)
}

func (rc *RouteConfig) FamilyMemberRoutes() {
	familyMember := rc.App.Group("/api/v1/family-members", rc.JWT.JWTProtected())
	familyMember.Post("/", rc.FamilyMemberController.Create)
	familyMember.Get("/:id", rc.FamilyMemberController.GetById)
	familyMember.Get("/", rc.FamilyMemberController.GetAll)
	familyMember.Put("/:id", rc.FamilyMemberController.Update)
	familyMember.Delete("/:id", rc.FamilyMemberController.Delete)
}

func (rc *RouteConfig) ClaimRoutes() {
	claim := rc.App.Group("/api/v1/claims", rc.JWT.JWTProtected())
	claim.Post("/", rc.ClaimController.CreateClaim)
	claim.Get("/get-patients", rc.ClaimController.GetAllPatient)
	claim.Get("/get-benefits/:patientId", rc.ClaimController.GetAllBenefits)
	claim.Put("/:id", rc.ClaimController.Update)
}
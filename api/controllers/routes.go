package controllers

import (
	"os"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
)

// InitializeRoutes initialize routing
func (s *Server) InitializeRoutes(rg *routing.RouteGroup) {

	// Home Routing
	rg.Any("/", s.Home())

	// Login Routing
	rg.Post("/login", s.Login())
	rg.Put("/meter/updatebynumber/<id>", s.UpdateMeterByMeterNumber())

	rg.Use(auth.JWT(os.Getenv("API_SECRET")))

	rg.Get("/user-stats", s.UserStatsHandler())
	rg.Get("/meter-stats", s.MeterStatsHandler())

	// Companies Route
	rg.Get("/companies/get/<id>", s.GetCompaniesHandler())
	rg.Post("/companies/create", s.CreateCompaniesHandler())
	rg.Get("/companies/list", s.ListCompaniesHandler())
	rg.Get("/companies/count", s.CountCompaniesHandler())
	rg.Put("/companies/update/<id>", s.UpdateCompaniesHandler())
	rg.Delete("/companies/delete/<id>", s.DeleteCompaniesHandler())

	// Users routings
	rg.Post("/users", s.CreateUserController())
	rg.Get("/users", s.ListUsersController())
	rg.Get("/users/<id>", s.GetUserController())
	rg.Put("/users", s.UpdateUserController())
	rg.Delete("/users/<id>", s.DeleteUserController())
	rg.Get("/user/count", s.CountUsers())

	// Meters Routing
	rg.Post("/meters", s.CreateMeter())
	rg.Get("/meters", s.ListMeters())
	rg.Get("/meter/get/<id>", s.GetMeterByID())
	rg.Get("/meter/getbynumber/<id>", s.GetMeterByMeterNumber())
	// rg.Get("/meter", s.GetVMeterDetailsByRegNoController())
	rg.Put("/meter/<id>", s.UpdateMeter())
	rg.Put("/meter/updatebymeterno/<id>", s.UpdateMeterByMeterNumber())
	rg.Delete("/meter/<id>", s.DeleteMeter())
	rg.Get("/meter/telemetry/<id>", s.GetMeterTelemetryController())
	rg.Get("/meter/count", s.CountMeters())

	// Gateways Routing
	rg.Post("/gateways", s.CreateGatewaysHandler())
	rg.Get("/gateways", s.ListGatewaysHandler())
	rg.Get("/gateway/get/<id>", s.GetGatewayHandler())
	rg.Put("/gateway/<id>", s.UpdateGatewayHandler())
	rg.Delete("/gateway/<id>", s.DeleteGatewayHandler())
	rg.Get("/gateway/count", s.CountGateways())

	// Settings Route
	rg.Post("/settings/permissions", s.CreatePermissions())
	rg.Get("/settings/permissions", s.ListPermissions())
	rg.Get("/settings/roles", s.ListRoles())

	rg.Get("/pricing/plan", s.ListPricePlanHandler())
	rg.Post("/pricing/plan", s.CreatePricePlanHandler())
	rg.Delete("/pricing/plan/<id>", s.DeletePricePlanHandler())

	rg.Get("/invoices/list", s.ListPaymentsHandler())
	rg.Post("/invoice/calculate", s.CalculateInvoiceHandler())
	rg.Delete("/invoice/pay", s.DeletePricePlanHandler())

	rg.Get("/building/get/<id>", s.GetBuildingHandler())
	rg.Get("/building/get-with-houses/<id>", s.GetBuildingHouseNumbersHandler())
	rg.Get("/buildings/count", s.CountBuildingsHandler())
	rg.Get("/buildings/list", s.ListBuildingsHandler())
	rg.Post("/buildings/create", s.CreateBuildingHandler())
	rg.Put("/buildings/update/<id>", s.UpdateBuildingHandler())
	rg.Delete("/buildings/delete/<id>", s.DeleteBuildingHandler())
	rg.Post("/buildings/house-number", s.CreateHouseNumbersHandler())
}

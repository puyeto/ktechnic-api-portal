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

	rg.Use(auth.JWT(os.Getenv("API_SECRET")))

	rg.Get("/user-stats", s.UserStatsHandler())
	rg.Get("/meter-stats", s.MeterStatsHandler())

	// Companies Route
	rg.Post("/companies", s.CreateCompanies())
	rg.Get("/companies", s.ListCompanies())

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
	rg.Get("/meter/get/<id>", s.GetMeter())
	// rg.Get("/meter", s.GetVMeterDetailsByRegNoController())
	rg.Put("/meter/<id>", s.UpdateMeter())
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
}

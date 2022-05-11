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

	rg.Get("/dashboard-stats", s.DashboardStatsHandler())

	// Companies Route
	rg.Post("/companies", s.CreateCompanies())
	rg.Get("/companies", s.ListCompanies())

	// Users routings
	rg.Post("/users", s.CreateUserController())
	rg.Get("/users", s.ListUsersController())
	rg.Get("/users/<id>", s.GetUserController())
	rg.Put("/users", s.UpdateUserController())
	rg.Delete("/users/<id>", s.DeleteUserController())
	rg.Get("/user/drivers", s.GetDriversController())

	// Meters Routing
	rg.Post("/meters", s.CreateMeter())
	rg.Get("/meters", s.ListMeters())
	rg.Get("/meter/get/<id>", s.GetMeter())
	// rg.Get("/meter", s.GetVMeterDetailsByRegNoController())
	rg.Put("/meters", s.UpdateMeter())
	rg.Delete("/meters/<id>", s.DeleteMeter())
	rg.Get("/meter/telemetry/<id>", s.GetMeterTelemetryController())

	// Gateways Routing
	rg.Post("/gateways", s.CreateGatewaysHandler())
	rg.Get("/gateways", s.ListGatewaysHandler())
	rg.Get("/gateway/get/<id>", s.GetGatewayHandler())
	rg.Put("/gateways", s.UpdateGatewayHandler())
	rg.Delete("/gateway/<id>", s.DeleteGatewayHandler())

	// Vehicle Routing
	rg.Post("/vehicles", s.CreateVehicle())
	rg.Get("/vehicles", s.ListVehicles())
	rg.Get("/vehicle/<id>", s.GetVehicle())
	rg.Get("/vehicle", s.GetVehicleDetailsByRegNoController())
	rg.Put("/vehicles", s.UpdateVehicle())
	rg.Delete("/vehicles/<id>", s.DeleteVehicle())

	rg.Post("/refuels/vehicle", s.RefuelVehicleController())
	rg.Get("/refuels/vehicle/<id>", s.GetRefuelsByVehicleIDController())

	rg.Post("/odometer/vehicle", s.AddVehicleOdometerReadingController())
	rg.Get("/odometer-logs/vehicle/<id>", s.GetVehicleOdometerReadingByVehicleIDController())

	// Vehicle Route Routing
	rg.Post("/routes", s.CreateVehicleRoute())
	rg.Get("/routes", s.ListVehicleRoutes())

	// Settings Route
	rg.Post("/settings/permissions", s.CreatePermissions())
	rg.Get("/settings/permissions", s.ListPermissions())
	rg.Get("/settings/roles", s.ListRoles())
	rg.Get("/settings/service-types", s.ListServiceTypesController())
	rg.Post("/settings/service-types", s.CreateServiceTypesController())
}

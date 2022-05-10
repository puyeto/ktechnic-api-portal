package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// Home ...
func (server *Server) Home() routing.Handler {
	return func(c *routing.Context) error {
		return c.Write(map[string]interface{}{
			"message": "Welcome To Lectrotel Energy Monitoring System",
		})
	}
}

// DashboardStatsHandler ...
func (server *Server) DashboardStatsHandler() routing.Handler {
	return func(c *routing.Context) error {
		ds := models.DashboardStats{}

		ds.CompanyID = auth.ExtractCompanyID(c)

		return c.Write(map[string]interface{}{
			"response": ds.GetUserDashboardStats(server.DB),
		})
	}
}

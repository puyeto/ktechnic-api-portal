package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// CreateUser ...
func (server *Server) CreateVehicleRoute() routing.Handler {
	return func(c *routing.Context) error {
		var route models.VehicleRoutes
		if err := c.Read(&route); err != nil {
			return errors.BadRequest(err.Error())
		}
		route.Prepare()

		err := route.Validate("")
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		route.UpdatedBy = auth.ExtractTokenID(c)
		route.CompanyID = auth.ExtractCompanyID(c)

		userCreated, err := route.SaveVehicleRoute(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}
		return c.Write(map[string]interface{}{
			"response": userCreated,
		})
	}
}

// ListVehicleRoutes ...
func (server *Server) ListVehicleRoutes() routing.Handler {
	return func(c *routing.Context) error {
		route := models.VehicleRoutes{}

		route.CompanyID = auth.ExtractCompanyID(c)

		routes, err := route.ListVehicleRoutes(server.DB)
		if err != nil {
			return errors.NoContentFound(err.Error())
		}
		return c.Write(map[string]interface{}{
			"response": routes,
		})
	}
}

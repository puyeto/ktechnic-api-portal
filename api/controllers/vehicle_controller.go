package controllers

import (
	"regexp"
	"strconv"
	"strings"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// CreateVehicle ..
func (server *Server) CreateVehicle() routing.Handler {
	return func(c *routing.Context) error {
		var vehicle models.Vehicle
		if err := c.Read(&vehicle); err != nil {
			return errors.BadRequest(err.Error())
		}

		vehicle.Prepare()
		err := vehicle.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		vehicle.AddedBy = auth.ExtractTokenID(c)
		vehicle.CompanyID = auth.ExtractCompanyID(c)

		vehicleCreated, err := vehicle.SaveVehicle(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": vehicleCreated,
		})
	}
}

// ListVehicles ...
func (server *Server) ListVehicles() routing.Handler {
	return func(c *routing.Context) error {
		vehicle := models.Vehicle{}

		vehicle.CompanyID = auth.ExtractCompanyID(c)

		vehicles, err := vehicle.ListAllVehicles(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		// responses.JSON(w, http.StatusOK, vehicles)
		return c.Write(map[string]interface{}{
			"response": vehicles,
		})
	}
}

// GetVehicle ...
func (server *Server) GetVehicle() routing.Handler {
	return func(c *routing.Context) error {
		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		vehicle := models.Vehicle{}
		vehicleReceived, err := vehicle.FindVehicleByID(server.DB, uint64(vid))
		if err != nil {
			return errors.NoContentFound(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": vehicleReceived,
		})
	}
}

// GetVehicleDetailsBy ...
func (server *Server) GetVehicleDetailsByRegNoController() routing.Handler {
	return func(c *routing.Context) error {
		vehicle := models.Vehicle{}
		reg := c.Query("reg", "")

		if reg != "" {
			rec, _ := regexp.Compile("[^a-zA-Z0-9]+")
			vehicle.VehicleStringID = strings.ToLower(rec.ReplaceAllString(reg, "_"))
		}

		vehicleReceived, err := vehicle.FindVehicleByQueryString(server.DB)
		if err != nil {
			return errors.NoContentFound(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": vehicleReceived,
		})
	}
}

// UpdateVehicle ...
func (server *Server) UpdateVehicle() routing.Handler {
	return func(c *routing.Context) error {
		var vehicle models.Vehicle
		if err := c.Read(&vehicle); err != nil {
			return errors.BadRequest(err.Error())
		}
		vehicle.Prepare()

		if err := vehicle.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		vehicleUpdated, err := vehicle.UpdateAVehicle(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": vehicleUpdated,
		})
	}
}

// DeleteVehicle ...
func (server *Server) DeleteVehicle() routing.Handler {
	return func(c *routing.Context) error {
		var vehicle models.Vehicle

		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		_, err = vehicle.DeleteAVehicle(server.DB, uint32(vid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": "success",
		})
	}
}

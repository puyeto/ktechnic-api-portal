package controllers

import (
	"fmt"
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// GetVehicleOdometerReadingByVehicleIDController ...
func (server *Server) GetVehicleOdometerReadingByVehicleIDController() routing.Handler {
	return func(c *routing.Context) error {
		var odometer = models.Odometers{}

		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		serv, err := odometer.GetOdometerByVehicleID(server.DB, uint64(vid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		// responses.JSON(w, http.StatusOK, vehicles)
		return c.Write(map[string]interface{}{
			"response": serv,
		})
	}
}

// AddVehicleOdometerReadingController ...
func (server *Server) AddVehicleOdometerReadingController() routing.Handler {
	return func(c *routing.Context) error {
		var odometer models.Odometers
		if err := c.Read(&odometer); err != nil {
			fmt.Println(err)
			return errors.BadRequest(err.Error())
		}
		odometer.Prepare()

		err := odometer.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		odometer.AddedBy = auth.ExtractTokenID(c)

		odometerCreated, err := odometer.Save(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": odometerCreated,
		})
	}
}

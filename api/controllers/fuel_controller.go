package controllers

import (
	"fmt"
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// GetRefuelsByVehicleIDController ...
func (server *Server) GetRefuelsByVehicleIDController() routing.Handler {
	return func(c *routing.Context) error {
		var fuel = models.Fuel{}

		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		serv, err := fuel.GetRefuelsByVehicleID(server.DB, uint64(vid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		// responses.JSON(w, http.StatusOK, vehicles)
		return c.Write(map[string]interface{}{
			"response": serv,
		})
	}
}

// RefuelVehicleController ...
func (server *Server) RefuelVehicleController() routing.Handler {
	return func(c *routing.Context) error {
		var fuel models.Fuel
		if err := c.Read(&fuel); err != nil {
			fmt.Println(err)
			return errors.BadRequest(err.Error())
		}
		fuel.Prepare()

		err := fuel.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		fuel.AddedBy = auth.ExtractTokenID(c)

		// Get Odometer id after save
		if fuel.OdometerValue > 0 {
			odometer := models.InitializeOdometer(fuel.VehicleID, fuel.OdometerValue, fuel.DateFueled, fuel.AddedBy)
			odometer.Prepare()
			odometerCreated, err := odometer.Save(server.DB)
			if err != nil {
				return errors.InternalServerError(err.Error())
			}
			fuel.OdometerID = odometerCreated.ID
		}

		// Get Expense id after save
		expense := models.InitializeExpenses(fuel.VehicleID, 73, fuel.TotalPrice, fuel.Description, fuel.DateFueled, fuel.AddedBy)
		expense.Prepare()
		expenseCreated, err := expense.Save(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}
		fuel.ExpenseID = expenseCreated.ID

		fuelCreated, err := fuel.Save(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": fuelCreated,
		})
	}
}

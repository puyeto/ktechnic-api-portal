package controllers

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

func (server *Server) ListPricePlanHandler() routing.Handler {
	return func(c *routing.Context) error {
		plan := models.PricePlan{}

		plan.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		plan.AddedBy = auth.ExtractTokenID(c)
		plans, err := plan.ListPricePlan(server.DB, roleid)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(plans)
	}
}

func (server *Server) CreatePricePlanHandler() routing.Handler {
	return func(c *routing.Context) error {
		var pplan models.PricePlan
		if err := c.Read(&pplan); err != nil {
			return errors.BadRequest(err.Error())
		}

		err := pplan.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		pplan.AddedBy = auth.ExtractTokenID(c)
		pplan.CompanyID = auth.ExtractCompanyID(c)

		planCreated, err := pplan.SavePricePlan(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(planCreated)
	}
}

// Delete Price Plan ...
func (server *Server) DeletePricePlanHandler() routing.Handler {
	return func(c *routing.Context) error {
		var plan models.PricePlan

		pid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		err = plan.DeletePricePlan(server.DB, uint32(pid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": "success",
		})
	}
}

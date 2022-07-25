package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

func (server *Server) ListPaymentsHandler() routing.Handler {
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

func (server *Server) CalculateInvoiceHandler() routing.Handler {
	return func(c *routing.Context) error {
		var calc models.CalculateTotalsDetails
		if err := c.Read(&calc); err != nil {
			return errors.BadRequest(err.Error())
		}

		err := calc.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		metercalc, err := calc.CalculateTotal(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(metercalc)
	}
}

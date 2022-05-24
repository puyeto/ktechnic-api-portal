package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// CreateContributionFields ...
func (server *Server) CreateCompanies() routing.Handler {
	return func(c *routing.Context) error {
		var company models.Companies
		if err := c.Read(&company); err != nil {
			return errors.BadRequest(err.Error())
		}
		company.Prepare()

		err := company.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		uid := auth.ExtractTokenID(c)
		company.AddedBy = uid

		companyCreated, err := company.SaveCompanyDetails(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(companyCreated)
	}
}

func (server *Server) ListCompanies() routing.Handler {
	return func(c *routing.Context) error {
		var company = models.Companies{}

		companies, err := company.List(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(companies)
	}
}

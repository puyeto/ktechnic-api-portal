package controllers

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// GetCompaniesHandler ..
func (server *Server) GetCompaniesHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			id      = stringToUInt64(c.Param("id"))
			company = models.Companies{}
		)
		if id == 0 {
			return errors.InternalServerError("Invalid Company Detail")
		}

		company.ID = uint32(id)
		res, err := company.Get(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(res)
	}
}

// CreateContributionFields ...
func (server *Server) CreateCompaniesHandler() routing.Handler {
	return func(c *routing.Context) error {
		var company models.Companies
		if err := c.Read(&company); err != nil {
			return errors.BadRequest(err.Error())
		}

		company.Prepare()
		if err := company.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		company.AddedBy = auth.ExtractTokenID(c)
		companyCreated, err := company.Create(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(companyCreated)
	}
}

// Count Meters ...
func (server *Server) CountCompaniesHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			company = models.Companies{}
			roleid  = auth.ExtractRoleID(c)
		)

		company.AddedBy = auth.ExtractTokenID(c)
		count := company.Count(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		return c.Write(map[string]int{
			"result": count,
		})
	}
}

func (server *Server) ListCompaniesHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			company = models.Companies{}
			page    = parseInt(c.Query("page"), 1)
			perPage = parseInt(c.Query("per_page"), 0)
			roleid  = auth.ExtractRoleID(c)
		)
		company.AddedBy = auth.ExtractTokenID(c)

		count := company.Count(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		paginatedList := getPaginatedListFromRequest(c, count, page, perPage)
		res, err := company.List(server.DB, roleid, paginatedList.Offset(), paginatedList.Limit())
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		paginatedList.Items = res
		return c.Write(paginatedList)
	}
}

// UpdateCompaniesHandler ...
func (server *Server) UpdateCompaniesHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			company models.Companies
			id      = stringToUInt64(c.Param("id"))
		)
		if id == 0 {
			return errors.InternalServerError("Invalid Company Detail")
		}

		company.ID = uint32(id)
		received, err := company.Get(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		if err := c.Read(&received); err != nil {
			return errors.BadRequest(err.Error())
		}

		if err := received.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		meterUpdated, err := received.Update(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meterUpdated)
	}
}

// Delete ...
func (server *Server) DeleteCompaniesHandler() routing.Handler {
	return func(c *routing.Context) error {
		var company models.Companies

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		company.ID = uint32(id)
		if err := company.Delete(server.DB); err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"result": "success",
		})
	}
}

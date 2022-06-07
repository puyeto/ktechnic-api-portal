package controllers

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// Get ...
func (server *Server) GetBuildingHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			id    = stringToUInt64(c.Param("id"))
			build = models.Buildings{}
		)
		if id == 0 {
			return errors.InternalServerError("Invalid Building Detail")
		}

		build.ID = uint32(id)
		res, err := build.Get(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(res)
	}
}

// Get Building House Numbers ...
func (server *Server) GetBuildingHouseNumbersHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			id    = stringToUInt64(c.Param("id"))
			build = models.BuildingWithHouseNo{}
		)
		if id == 0 {
			return errors.InternalServerError("Invalid Building Detail")
		}

		build.ID = uint32(id)
		res, err := build.GetWithHouseNumbers(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(res)
	}
}

// Count Meters ...
func (server *Server) CountBuildingsHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			build  = models.Buildings{}
			roleid = auth.ExtractRoleID(c)
		)
		build.CompanyID = auth.ExtractCompanyID(c)
		build.AddedBy = auth.ExtractTokenID(c)

		count := build.Count(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		return c.Write(map[string]int{
			"result": count,
		})
	}
}

func (server *Server) ListBuildingsHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			build   = models.Buildings{}
			page    = parseInt(c.Query("page"), 1)
			perPage = parseInt(c.Query("per_page"), 0)
			roleid  = auth.ExtractRoleID(c)
		)

		build.CompanyID = auth.ExtractCompanyID(c)
		build.AddedBy = auth.ExtractTokenID(c)

		count := build.Count(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		paginatedList := getPaginatedListFromRequest(c, count, page, perPage)
		res, err := build.List(server.DB, roleid, paginatedList.Offset(), paginatedList.Limit())
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		paginatedList.Items = res
		return c.Write(paginatedList)
	}
}

// CreateContributionFields ...
func (server *Server) CreateBuildingHandler() routing.Handler {
	return func(c *routing.Context) error {
		var build models.Buildings
		if err := c.Read(&build); err != nil {
			return errors.BadRequest(err.Error())
		}

		build.Prepare()
		if err := build.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		build.AddedBy = auth.ExtractTokenID(c)
		build.CompanyID = auth.ExtractCompanyID(c)
		res, err := build.Create(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(res)
	}
}

// UpdateBuildingHandler
func (server *Server) UpdateBuildingHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			build models.Buildings
			id    = stringToUInt64(c.Param("id"))
		)
		if id == 0 {
			return errors.InternalServerError("Invalid Building Detail")
		}

		build.ID = uint32(id)
		received, err := build.Get(server.DB)
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
func (server *Server) DeleteBuildingHandler() routing.Handler {
	return func(c *routing.Context) error {
		var build models.Buildings

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		build.ID = uint32(id)
		if err := build.Delete(server.DB); err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"result": "success",
		})
	}
}

// Create Building House Numbers
func (server *Server) CreateHouseNumbersHandler() routing.Handler {
	return func(c *routing.Context) error {
		var (
			build models.BuildingHouseNumber
		)

		if err := c.Read(&build); err != nil {
			return errors.BadRequest(err.Error())
		}

		if err := build.Create(server.DB); err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"result": "success",
		})
	}
}

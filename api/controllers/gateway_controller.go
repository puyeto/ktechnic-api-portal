package controllers

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// CreateGatewaysHandler ..
func (server *Server) CreateGatewaysHandler() routing.Handler {
	return func(c *routing.Context) error {
		var gateway models.Gateway
		if err := c.Read(&gateway); err != nil {
			return errors.BadRequest(err.Error())
		}

		gateway.Prepare()
		err := gateway.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		gateway.AddedBy = auth.ExtractTokenID(c)
		gateway.CompanyID = auth.ExtractCompanyID(c)

		gatewayCreated, err := gateway.SaveGateway(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(gatewayCreated)
	}
}

// Count Gateways ...
func (server *Server) CountGateways() routing.Handler {
	return func(c *routing.Context) error {
		gate := models.Gateway{}
		gate.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		gate.AddedBy = auth.ExtractTokenID(c)

		count := gate.CountGateways(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		return c.Write(map[string]int{
			"result": count,
		})
	}
}

// ListGatewaysHandler ...
func (server *Server) ListGatewaysHandler() routing.Handler {
	return func(c *routing.Context) error {
		page := parseInt(c.Query("page"), 1)
		perPage := parseInt(c.Query("per_page"), 0)
		gateway := models.Gateway{}

		gateway.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		gateway.AddedBy = auth.ExtractTokenID(c)

		gateways, err := gateway.ListAllGateways(server.DB, roleid, uint32(page), uint32(perPage))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		// responses.JSON(w, http.StatusOK, gateways)
		return c.Write(gateways)
	}
}

// GetGatewayHandler ...
func (server *Server) GetGatewayHandler() routing.Handler {
	return func(c *routing.Context) error {
		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		gateway := models.Gateway{}
		gatewayReceived, err := gateway.FindGatewayByID(server.DB, uint64(vid))
		if err != nil {
			return errors.NoContentFound(err.Error())
		}

		return c.Write(gatewayReceived)
	}
}

// UpdateGatewayHandler ...
func (server *Server) UpdateGatewayHandler() routing.Handler {
	return func(c *routing.Context) error {
		var gateway models.Gateway
		if err := c.Read(&gateway); err != nil {
			return errors.BadRequest(err.Error())
		}
		gateway.Prepare()

		if err := gateway.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		gatewayUpdated, err := gateway.UpdateAGateway(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(gatewayUpdated)
	}
}

// DeleteGatewayHandler ...
func (server *Server) DeleteGatewayHandler() routing.Handler {
	return func(c *routing.Context) error {
		var gateway models.Gateway

		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		err = gateway.DeleteAGateway(server.DB, uint32(vid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"result": "success",
		})
	}
}

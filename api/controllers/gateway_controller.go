package controllers

import (
	"strconv"
	"time"

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

		count := gateway.CountGateways(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		paginatedList := getPaginatedListFromRequest(c, count, page, perPage)
		gateways, err := gateway.ListAllGateways(server.DB, roleid, paginatedList.Offset(), paginatedList.Limit())
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		paginatedList.Items = gateways
		return c.Write(paginatedList)
	}
}

// GetGatewayHandler ...
func (server *Server) GetGatewayHandler() routing.Handler {
	return func(c *routing.Context) error {
		gid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		gateway := models.Gateway{}
		gateway.ID = gid
		gatewayReceived, err := gateway.FindGatewayByID(server.DB)
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
		gid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		gateway.ID = gid
		gatewayReceived, err := gateway.FindGatewayByID(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		if err := c.Read(&gatewayReceived); err != nil {
			return errors.BadRequest(err.Error())
		}

		gatewayReceived.UpdatedAt = time.Now()
		if err := gatewayReceived.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		gatewayUpdated, err := gatewayReceived.UpdateAGateway(server.DB)
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

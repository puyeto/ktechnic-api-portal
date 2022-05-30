package controllers

import (
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// Home ...
func (server *Server) Home() routing.Handler {
	return func(c *routing.Context) error {
		return c.Write(map[string]interface{}{
			"message": "Welcome To KTechnics Energy System",
		})
	}
}

// UserStatsHandler ...
func (server *Server) UserStatsHandler() routing.Handler {
	return func(c *routing.Context) error {
		userid := auth.ExtractTokenID(c)
		companyid := auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)

		// Get Meter Count
		meter := models.Meter{}
		meter.CompanyID = companyid
		meter.AddedBy = userid
		mCount := meter.CountMeters(server.DB, roleid)

		// Get User Count
		user := models.User{}
		user.CompanyID = companyid
		user.UpdatedBy = userid
		uCount := user.CountUsers(server.DB, roleid)

		// Get Gateway Count
		gate := models.Gateway{}
		gate.CompanyID = companyid
		gate.AddedBy = userid
		gCount := gate.CountGateways(server.DB, roleid)

		return c.Write(models.UserStats{
			MeterCount:   mCount,
			UserCount:    uCount,
			GatewayCount: gCount,
		})
	}
}

// MeterStatsHandler ...
func (server *Server) MeterStatsHandler() routing.Handler {
	return func(c *routing.Context) error {
		meter := models.Meter{}
		meter.ID = uint32(parseInt(c.Query("meter_id"), 0))

		// Check if meter exists
		meterDetails, err := meter.GetMeterByID(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		t := time.Now()
		return c.Write(models.MaterStats{
			MeterNumber:   meterDetails.MeterNumber,
			MeterLastSeen: t.Format("Mon Jan 02 15:04:05"),
			MeterDetails:  meterDetails,
		})
	}
}

package controllers

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/app"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// CreateMeter ..
func (server *Server) CreateMeter() routing.Handler {
	return func(c *routing.Context) error {
		var meter models.Meter
		if err := c.Read(&meter); err != nil {
			return errors.BadRequest(err.Error())
		}

		meter.Prepare()
		err := meter.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		meter.AddedBy = auth.ExtractTokenID(c)
		meter.CompanyID = auth.ExtractCompanyID(c)

		meterCreated, err := meter.SaveMeter(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meterCreated)
	}
}

// ListMeters ...
func (server *Server) ListMeters() routing.Handler {
	return func(c *routing.Context) error {
		meter := models.Meter{}

		meter.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		meter.AddedBy = auth.ExtractTokenID(c)
		meters, err := meter.ListAllMeters(server.DB, roleid)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meters)
	}
}

// GetMeter ...
func (server *Server) GetMeter() routing.Handler {
	return func(c *routing.Context) error {
		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		meter := models.Meter{}
		meterReceived, err := meter.FindMeterByID(server.DB, uint64(vid))
		if err != nil {
			return errors.NoContentFound(err.Error())
		}

		return c.Write(meterReceived)
	}
}

// UpdateMeter ...
func (server *Server) UpdateMeter() routing.Handler {
	return func(c *routing.Context) error {
		var meter models.Meter
		if err := c.Read(&meter); err != nil {
			return errors.BadRequest(err.Error())
		}
		meter.Prepare()

		if err := meter.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		meterUpdated, err := meter.UpdateAMeter(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": meterUpdated,
		})
	}
}

// DeleteMeter ...
func (server *Server) DeleteMeter() routing.Handler {
	return func(c *routing.Context) error {
		var meter models.Meter

		vid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		_, err = meter.DeleteAMeter(server.DB, uint32(vid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": "success",
		})
	}
}

// GetMeterTelemetryController ...
func (server *Server) GetMeterTelemetryController() routing.Handler {
	return func(c *routing.Context) error {
		order := c.Query("order_by", "desc")
		filterfrom, err := strconv.Atoi(c.Query("filter_from", "0"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		filterto, err := strconv.Atoi(c.Query("filter_to", "0"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		mid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		meter := models.Meter{}
		count := meter.CountMeterTelemetryByID(app.MongoDB, uint64(mid), uint64(filterfrom), uint64(filterto))
		var paginatedList *app.PaginatedList

		if count > 0 {
			paginatedList = getPaginatedListFromRequest(c, count)
			meterReceived, err := meter.FindMeterTelemetryByID(app.MongoDB, uint64(mid), order, paginatedList.Offset(), paginatedList.Limit(), uint64(filterfrom), uint64(filterto))
			if err != nil {
				return errors.NoContentFound(err.Error())
			}
			paginatedList.Items = meterReceived
		}

		return c.Write(paginatedList)
	}
}

package controllers

import (
	"strconv"
	"time"

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

// Count Meters ...
func (server *Server) CountMeters() routing.Handler {
	return func(c *routing.Context) error {
		meter := models.Meter{}
		meter.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		meter.AddedBy = auth.ExtractTokenID(c)

		count := meter.CountMeters(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		return c.Write(map[string]int{
			"result": count,
		})
	}
}

// ListMeters ...
func (server *Server) ListMeters() routing.Handler {
	return func(c *routing.Context) error {
		page := parseInt(c.Query("page"), 1)
		perPage := parseInt(c.Query("per_page"), 0)
		meter := models.Meter{}

		meter.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		meter.AddedBy = auth.ExtractTokenID(c)

		count := meter.CountMeters(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		paginatedList := getPaginatedListFromRequest(c, count, page, perPage)
		meters, err := meter.ListAllMeters(server.DB, roleid, paginatedList.Offset(), paginatedList.Limit())
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		paginatedList.Items = meters
		return c.Write(paginatedList)
	}
}

// GetMeterByID ...
func (server *Server) GetMeterByID() routing.Handler {
	return func(c *routing.Context) error {
		mid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		var meter models.Meter
		meter.ID = uint32(mid)
		meterReceived, err := meter.GetMeterByID(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meterReceived)
	}
}

// GetMeterByMeterNumber
func (server *Server) GetMeterByMeterNumber() routing.Handler {
	return func(c *routing.Context) error {
		mno := stringToUInt64(c.Param("id"))
		if mno == 0 {
			return errors.InternalServerError("Invalid Meter Number")
		}

		var meter models.Meter
		meter.MeterNumber = mno
		meterReceived, err := meter.GetMeterByMeterNumber(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meterReceived)
	}
}

// UpdateMeter ...
func (server *Server) UpdateMeter() routing.Handler {
	return func(c *routing.Context) error {
		var meter models.Meter
		mid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		meter.ID = uint32(mid)
		meterReceived, err := meter.GetMeterByID(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		if err := c.Read(&meterReceived); err != nil {
			return errors.BadRequest(err.Error())
		}

		meterReceived.UpdatedAt = time.Now()
		if err := meterReceived.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		meterUpdated, err := meterReceived.UpdateAMeter(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meterUpdated)
	}
}

// UpdateMeterByMeterNumber
func (server *Server) UpdateMeterByMeterNumber() routing.Handler {
	return func(c *routing.Context) error {
		var meter models.Meter
		mno := stringToUInt64(c.Param("id"))
		if mno == 0 {
			return errors.InternalServerError("Invalid Meter Number")
		}

		meter.MeterNumber = mno
		meterReceived, err := meter.GetMeterByMeterNumber(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		if err := c.Read(&meterReceived); err != nil {
			return errors.BadRequest(err.Error())
		}

		meterReceived.UpdatedAt = time.Now()
		if err := meterReceived.Validate(); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		meterUpdated, err := meterReceived.UpdateAMeter(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(meterUpdated)
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

		if err := meter.DeleteAMeter(server.DB, uint32(vid)); err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"result": "success",
		})
	}
}

// GetMeterTelemetryController ...
func (server *Server) GetMeterTelemetryController() routing.Handler {
	return func(c *routing.Context) error {
		order := c.Query("order_by", "desc")
		page := parseInt(c.Query("page"), 1)
		perPage := parseInt(c.Query("per_page"), 0)
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
		paginatedList := getPaginatedListFromRequest(c, count, page, perPage)
		if count > 0 {
			meterReceived, err := meter.FindMeterTelemetryByID(app.MongoDB, uint64(mid), order, paginatedList.Offset(), paginatedList.Limit(), uint64(filterfrom), uint64(filterto))
			if err != nil {
				return errors.NoContentFound(err.Error())
			}
			paginatedList.Items = meterReceived
		}

		return c.Write(paginatedList)
	}
}

package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// ListServiceTypesController ...
func (server *Server) ListServiceTypesController() routing.Handler {
	return func(c *routing.Context) error {
		var ser = models.ServiceTypes{}

		serv, err := ser.List(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		// responses.JSON(w, http.StatusOK, vehicles)
		return c.Write(map[string]interface{}{
			"response": serv,
		})
	}
}

// CreateServiceTypesController ...
func (server *Server) CreateServiceTypesController() routing.Handler {
	return func(c *routing.Context) error {
		var ser models.ServiceTypes
		if err := c.Read(&ser); err != nil {
			return errors.BadRequest(err.Error())
		}
		ser.Prepare()

		err := ser.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		uid := auth.ExtractTokenID(c)
		ser.AddedBy = uid

		serCreated, err := ser.Save(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"response": serCreated,
		})
	}
}

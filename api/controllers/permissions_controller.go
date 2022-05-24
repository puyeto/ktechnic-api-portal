package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
	"github.com/ktechnics/ktechnics-api/api/utils/formaterror"
)

// ListPermissions ...
func (server *Server) ListPermissions() routing.Handler {
	return func(c *routing.Context) error {
		var perm = models.Permissions{}

		perms, err := perm.List(server.DB)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			return errors.InternalServerError(formattedError.Error())
		}

		// responses.JSON(w, http.StatusOK, vehicles)
		return c.Write(perms)
	}
}

// CreateContributionFields ...
func (server *Server) CreatePermissions() routing.Handler {
	return func(c *routing.Context) error {
		var perm models.Permissions
		if err := c.Read(&perm); err != nil {
			return errors.BadRequest(err.Error())
		}
		perm.Prepare()

		err := perm.Validate()
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		uid := auth.ExtractTokenID(c)
		perm.AddedBy = uid

		permCreated, err := perm.Save(server.DB)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(permCreated)
	}
}

// ListRoles ...
func (server *Server) ListRoles() routing.Handler {
	return func(c *routing.Context) error {
		var role = models.Roles{}

		// get company_id and role_id
		cid := auth.ExtractCompanyID(c)
		rid := auth.ExtractRoleID(c)

		roles, err := role.List(server.DB, cid, rid)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		// responses.JSON(w, http.StatusOK, vehicles)
		return c.Write(roles)
	}
}

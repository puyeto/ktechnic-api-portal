package controllers

import (
	"fmt"
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
	"github.com/ktechnics/ktechnics-api/api/utils/formaterror"
)

// CreateUserController ...
func (server *Server) CreateUserController() routing.Handler {
	return func(c *routing.Context) error {
		var user models.User
		if err := c.Read(&user); err != nil {
			fmt.Println(err)
			return errors.BadRequest(err.Error())
		}
		user.Prepare()

		err := user.Validate("")
		if err != nil {
			return errors.ValidationRequest(err.Error())
		}

		user.UpdatedBy = auth.ExtractTokenID(c)
		if user.CompanyID == 0 {
			user.CompanyID = auth.ExtractCompanyID(c)
		}

		userCreated, err := user.SaveUser(server.DB)
		if err != nil {
			fmt.Println(err)
			//  formaterror.FormatError(err.Error())
			return errors.InternalServerError(err.Error())
		}
		return c.Write(userCreated)
	}
}

// Count Users ...
func (server *Server) CountUsers() routing.Handler {
	return func(c *routing.Context) error {
		user := models.User{}
		user.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		user.UpdatedBy = auth.ExtractTokenID(c)

		count := user.CountUsers(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		return c.Write(map[string]int{
			"result": count,
		})
	}
}

// ListUsersController ...
func (server *Server) ListUsersController() routing.Handler {
	return func(c *routing.Context) error {
		page := parseInt(c.Query("page"), 1)
		perPage := parseInt(c.Query("per_page"), 0)
		user := models.User{}

		// get cid
		user.CompanyID = auth.ExtractCompanyID(c)
		roleid := auth.ExtractRoleID(c)
		user.UpdatedBy = auth.ExtractTokenID(c)

		count := user.CountUsers(server.DB, roleid)
		if count == 0 {
			return errors.InternalServerError("No Data Found")
		}

		paginatedList := getPaginatedListFromRequest(c, count, page, perPage)
		users, err := user.ListUsers(server.DB, roleid, paginatedList.Offset(), paginatedList.Limit())
		if err != nil {
			return errors.NoContentFound(err.Error())
		}

		paginatedList.Items = users
		return c.Write(paginatedList)
	}
}

// GetUserController ...
func (server *Server) GetUserController() routing.Handler {
	return func(c *routing.Context) error {
		uid, err := strconv.ParseUint(c.Query("id", "0"), 10, 32)
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		user := models.User{}
		userDetails, err := user.FindUserByID(server.DB, uint32(uid))
		if err != nil {
			return errors.NoContentFound(err.Error())
		}

		return c.Write(userDetails)
	}
}

// UpdateUserController ...
func (server *Server) UpdateUserController() routing.Handler {
	return func(c *routing.Context) error {
		var user models.User
		if err := c.Read(&user); err != nil {
			return errors.BadRequest(err.Error())
		}
		user.Prepare()

		user.UpdatedBy = auth.ExtractTokenID(c)

		if err := user.Validate("update"); err != nil {
			return errors.ValidationRequest(err.Error())
		}

		updatedUser, err := user.UpdateAUser(server.DB)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			return errors.InternalServerError(formattedError.Error())
		}

		return c.Write(updatedUser)
	}
}

// DeleteUserController ...
func (server *Server) DeleteUserController() routing.Handler {
	return func(c *routing.Context) error {
		user := models.User{}

		uid, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		_, err = user.DeleteAUser(server.DB, uint32(uid))
		if err != nil {
			return errors.InternalServerError(err.Error())
		}

		return c.Write(map[string]interface{}{
			"result": "success",
		})
	}
}

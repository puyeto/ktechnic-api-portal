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
		return c.Write(map[string]interface{}{
			"response": userCreated,
		})
	}
}

// ListUsersController ...
func (server *Server) ListUsersController() routing.Handler {
	return func(c *routing.Context) error {
		user := models.User{}

		// get cid
		user.CompanyID = auth.ExtractCompanyID(c)

		users, err := user.ListUsers(server.DB)
		if err != nil {
			return errors.NoContentFound(err.Error())
		}
		return c.Write(map[string]interface{}{
			"response": users,
		})
	}
}

// GetDriversController ...
func (server *Server) GetDriversController() routing.Handler {
	return func(c *routing.Context) error {
		user := models.User{}

		users, err := user.FindAllDrivers(server.DB)
		if err != nil {
			return errors.NoContentFound(err.Error())
		}
		return c.Write(map[string]interface{}{
			"response": users,
		})
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

		return c.Write(map[string]interface{}{
			"response": userDetails,
		})
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

		return c.Write(map[string]interface{}{
			"response": updatedUser,
		})
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
			"response": "success",
		})
	}
}

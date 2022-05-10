package controllers

import (
	"fmt"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/ktechnics/ktechnics-api/api/auth"
	"github.com/ktechnics/ktechnics-api/api/errors"
	"github.com/ktechnics/ktechnics-api/api/models"
	"github.com/ktechnics/ktechnics-api/api/services"
	"golang.org/x/crypto/bcrypt"
)

// Login ...
func (server *Server) Login() routing.Handler {
	return func(c *routing.Context) error {
		var user models.User
		if err := c.Read(&user); err != nil {
			return errors.Unauthorized(err.Error())
		}

		user.Prepare()
		err := user.Validate("login")
		if err != nil {
			return errors.Unauthorized(err.Error())
		}
		token, err := server.SignIn(user.Email, user.Password)
		if err != nil {
			return errors.InternalServerError(err.Error())
		}
		// responses.JSON(w, http.StatusOK, token)
		return c.Write(map[string]interface{}{
			"response": token,
		})
	}
}

// SignIn ...
func (server *Server) SignIn(email, password string) (m models.User, e error) {

	var err error

	user := models.User{}

	if err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error; err != nil {
		return m, err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return m, err
	}

	if user.CompanyID > 0 {
		if err := server.DB.Debug().Model(models.User{}).Where("id = ?", user.CompanyID).Take(&user.Company).Error; err != nil {
			return m, err
		}
	}

	token, err := auth.CreateToken(&user)
	user.Token = token

	// get permissions
	rows, err := server.DB.Debug().Table("user_permissions").Select("*").Where("user_id = ?", user.ID).Rows()
	if err != nil {
		fmt.Println(err)
		return m, err
	}

	finalRows, err := services.GetRowsColumsDetails(rows)
	if err != nil {
		return m, err
	}

	if len(finalRows) > 0 {
		user.PermissionField = finalRows
	}

	user.Password = ""
	if user.CompanyID == 0 {
		user.Company = models.Companies{}
	}

	return user, err
}

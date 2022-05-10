package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/ktechnics/ktechnics-api/api/models"
)

// CreateToken create new token
func CreateToken(u *models.User) (string, error) {
	token, err := auth.NewJWT(jwt.MapClaims{
		"id":         u.ID,
		"authorized": true,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
		"company_id": u.CompanyID,
		"role_id":    u.RoleID,
	}, os.Getenv("API_SECRET"))

	return token, err
}

// ExtractTokenID ...
func ExtractTokenID(c *routing.Context) uint32 {
	claims := c.Get("JWT").(*jwt.Token).Claims.(jwt.MapClaims)
	return uint32(claims["id"].(float64))
}

// ExtractCompanyID ...
func ExtractCompanyID(c *routing.Context) uint32 {
	claims := c.Get("JWT").(*jwt.Token).Claims.(jwt.MapClaims)
	return uint32(claims["company_id"].(float64))
}

func ExtractRoleID(c *routing.Context) uint32 {
	claims := c.Get("JWT").(*jwt.Token).Claims.(jwt.MapClaims)
	return uint32(claims["role_id"].(float64))
}

// TokenValid token validity
func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

// ExtractToken token
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

//Pretty display the claims licely in the terminal
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}

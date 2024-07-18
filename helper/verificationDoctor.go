package helper

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"medis/middleware"
	"medis/models"
	"strings"
)

func VerifyDoctorToken(db *gorm.DB, c echo.Context, secretKey []byte) (*models.Doctor, error) {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("Authorization token is missing")
	}

	authParts := strings.SplitN(tokenString, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return nil, errors.New("Invalid token format")
	}

	tokenString = authParts[1]

	username, err := middleware.VerifyToken(tokenString, secretKey)
	if err != nil {
		return nil, errors.New("Invalid token")
	}

	var doctor models.Doctor
	result := db.Where("username = ?", username).First(&doctor)
	if result.Error != nil {
		return nil, errors.New("Doctor not found")
	}

	return &doctor, nil
}

package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"html/template"
	"medis/models"
	"net/http"
)

func VerifyEmail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")

		var doctor models.Doctor
		result := db.Where("verification_token = ?", token).First(&doctor)
		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "Invalid verification token")
		}

		doctor.IsVerified = true
		doctor.VerificationToken = ""
		db.Save(&doctor)

		tmpl, err := template.ParseFiles("helper/verification.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return nil
	}
}

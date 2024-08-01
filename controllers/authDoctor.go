package controllers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"medis/auth"
	"medis/helper"
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

func RegisterDoctor(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		doctor := c.Get("doctor").(models.Doctor)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to hash password",
			})
		}

		doctor.Password = string(hashedPassword)
		doctor.Fullname = doctor.FirstName + " " + doctor.LastName
		db.Create(&doctor)

		doctor.Password = ""

		tokenString, err := auth.GenerateToken(doctor.Username, secretKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to generate token",
			})
		}

		response := map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Doctor account registered successfully",
			"token":   tokenString,
			"id":      doctor.ID,
		}

		return c.JSON(http.StatusOK, response)
	}
}

func SignInDoctor(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		existingDoctor := c.Get("doctor").(models.Doctor)

		// Generate Token
		tokenString, err := auth.GenerateToken(existingDoctor.Username, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to generate token",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Doctor login successful",
			"token":   tokenString,
			"id":      existingDoctor.ID})
	}
}

package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"medis/helper"
	"medis/models"
	"net/http"
	"regexp"
)

func ValidateDoctorRegistration(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var doctor models.Doctor
		if err := c.Bind(&doctor); err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		uniqueToken := helper.GenerateUniqueToken()
		doctor.VerificationToken = uniqueToken

		if len(doctor.FirstName) < 1 || len(doctor.FirstName) > 100 || !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(doctor.FirstName) {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "First Name must be between 1 and 100 characters and contain only letters",
			})
		}

		if len(doctor.LastName) > 0 {
			if len(doctor.LastName) > 100 || !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(doctor.LastName) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Last Name max 100 characters and contain only letters",
				})
			}
		}

		if len(doctor.Username) < 5 || len(doctor.Username) > 100 {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Username must be at least 5 characters and max 100 characters",
			})
		}

		if len(doctor.Password) < 8 || len(doctor.Password) > 100 || !helper.IsValidPassword(doctor.Password) {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Password must be at least 8 characters max 100 characters and contain a combination of letters and numbers",
			})
		}

		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		match, _ := regexp.MatchString(emailPattern, doctor.Email)
		if !match {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid email format",
			})
		}

		contactNumberRegex := regexp.MustCompile(`^\d{10,13}$`)
		if !contactNumberRegex.MatchString(doctor.ContactNumber) {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Contact number must be between 10 and 13 digits and contain only numbers",
			})
		}

		if err := helper.SendWelcomeEmail(doctor.Email, doctor.FirstName+" "+doctor.LastName, uniqueToken); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to send welcome email",
			})
		}

		c.Set("doctor", doctor)
		return next(c)
	}
}

func CheckDoctorUniqueness(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			doctor := c.Get("doctor").(models.Doctor)

			var existingDoctor models.Doctor
			result := db.Where("username = ?", doctor.Username).First(&existingDoctor)
			if result.Error == nil {
				return c.JSON(http.StatusConflict, helper.ErrorResponse{
					Code:    http.StatusConflict,
					Message: "Username already exists",
				})
			} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to check username",
				})
			}

			result = db.Where("email = ?", doctor.Email).First(&existingDoctor)
			if result.Error == nil {
				return c.JSON(http.StatusConflict, helper.ErrorResponse{
					Code:    http.StatusConflict,
					Message: "Email already exists",
				})
			} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Failed to check email",
				})
			}

			return next(c)
		}
	}
}

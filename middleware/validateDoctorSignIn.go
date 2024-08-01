package middleware

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"medis/helper"
	"medis/models"
	"net/http"
)

func ValidateDoctorSignIn(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var doctor models.Doctor
			if err := c.Bind(&doctor); err != nil {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				}
				return c.JSON(http.StatusBadRequest, errorResponse)
			}

			// Validasi Username
			if doctor.Username == "" {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Username is required",
				}
				return c.JSON(http.StatusBadRequest, errorResponse)
			}

			// Validasi Password
			if doctor.Password == "" {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Password is required",
				}
				return c.JSON(http.StatusBadRequest, errorResponse)
			}

			// Cek Username di Database
			var existingDoctor models.Doctor
			result := db.Where("username = ?", doctor.Username).First(&existingDoctor)
			if result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					errorResponse := helper.ErrorResponse{
						Code:    http.StatusUnauthorized,
						Message: "Invalid username",
					}
					return c.JSON(http.StatusUnauthorized, errorResponse)
				} else {
					errorResponse := helper.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: "Failed to check username",
					}
					return c.JSON(http.StatusInternalServerError, errorResponse)
				}
			}

			// Cek Password
			err := bcrypt.CompareHashAndPassword([]byte(existingDoctor.Password), []byte(doctor.Password))
			if err != nil {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Invalid password",
				}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			}

			// Cek Verifikasi Akun
			if !existingDoctor.IsVerified {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Account not verified. Please verify your account before logging in.",
				}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			}

			// Send Login Notification
			go func(email, username string) {
				if err := helper.SendLoginNotification(email, username); err != nil {
					fmt.Println("Failed to send notification email:", err)
				}
			}(existingDoctor.Email, existingDoctor.FirstName+" "+existingDoctor.LastName)

			c.Set("doctor", existingDoctor)
			return next(c)
		}
	}
}

func VerifyDoctorTokenMiddleware(db *gorm.DB, secretKey []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			doctor, err := helper.VerifyDoctorToken(db, c, secretKey)
			if err != nil {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: err.Error(),
				}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			}
			// Store the doctor object in the context
			c.Set("doctor", doctor)
			return next(c)
		}
	}
}

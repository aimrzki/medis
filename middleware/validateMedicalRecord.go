package middleware

import (
	"github.com/labstack/echo/v4"
	"medis/helper"
	"medis/models"
	"net/http"
)

func ValidateMedicalRecord(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var medicalRecord models.MedicalRecords
		if err := c.Bind(&medicalRecord); err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(medicalRecord.PatientName) < 1 || len(medicalRecord.PatientName) > 100 || !helper.ValidateLettersAndSpaces(medicalRecord.PatientName) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Patient name must be between 1 and 100 characters and contain only letters and spaces",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if !helper.ValidateDateFormat(medicalRecord.BirthDate) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Birth date must be in the format yyyy-mm-dd",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if !helper.ValidateEmailFormat(medicalRecord.Email) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid email format",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if !helper.ValidatePhoneNumber(medicalRecord.PhoneNumber) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Phone number must contain only digits and be at most 13 characters long",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(medicalRecord.Diagnosis) < 1 || len(medicalRecord.Diagnosis) > 3000 {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Diagnosis must be between 1 and 3000 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(medicalRecord.Prescription) < 1 || len(medicalRecord.Prescription) > 3000 {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Prescription must be between 1 and 3000 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(medicalRecord.CareSuggestion) < 1 || len(medicalRecord.CareSuggestion) > 3000 {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Care suggestion must be between 1 and 3000 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if err := helper.SendMedicalRecordNotification(medicalRecord.Email, medicalRecord.PatientName, medicalRecord.Diagnosis, medicalRecord.Prescription, medicalRecord.CareSuggestion); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to send medical record notification",
			})
		}

		// Store the medicalRecord object in the context
		c.Set("medicalRecord", medicalRecord)
		return next(c)
	}
}

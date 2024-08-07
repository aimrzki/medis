package controllers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"medis/helper"
	"medis/models"
	"net/http"
	"strconv"
	"time"
)

func AddMedicalRecord(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		doctor := c.Get("doctor").(models.Doctor)
		medicalRecord := c.Get("medicalRecord").(models.MedicalRecords)

		medicalRecord.DoctorID = doctor.ID

		if err := db.Create(&medicalRecord).Error; err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create medical record",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		successResponse := map[string]interface{}{
			"code":    http.StatusCreated,
			"error":   false,
			"message": "Medical record created successfully",
			"data":    medicalRecord,
		}
		return c.JSON(http.StatusCreated, successResponse)
	}
}

func GetMedicalRecordsByDoctor(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		doctor := c.Get("doctor").(*models.Doctor)

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		searching := c.QueryParam("searching")

		var medicalRecords []models.MedicalRecords
		query := db.Where("doctor_id = ?", doctor.ID).
			Offset(offset).
			Limit(limit).
			Order("id DESC")

		if searching != "" {
			searchPattern := "%" + searching + "%"
			query = query.Where(
				db.Where("patient_name ILIKE ?", searchPattern).
					Or("email ILIKE ?", searchPattern))
		}

		if err := query.Find(&medicalRecords).Error; err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to fetch medical records",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		var totalRecords int64
		countQuery := db.Model(&models.MedicalRecords{}).Where("doctor_id = ?", doctor.ID)
		if searching != "" {
			searchPattern := "%" + searching + "%"
			countQuery = countQuery.Where("patient_name ILIKE ?", searchPattern)
		}
		countQuery.Count(&totalRecords)

		response := map[string]interface{}{
			"code":         http.StatusOK,
			"error":        false,
			"message":      "Medical records fetched successfully",
			"data":         medicalRecords,
			"totalRecords": totalRecords,
			"page":         page,
			"limit":        limit,
		}
		return c.JSON(http.StatusOK, response)
	}
}

func GetMedicalRecordByID(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		doctor := c.Get("doctor").(*models.Doctor)

		recordID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid record ID",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		var medicalRecord models.MedicalRecords
		if err := db.Where("id = ? AND doctor_id = ?", recordID, doctor.ID).First(&medicalRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorResponse := helper.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Medical record not found or access denied",
				}
				return c.JSON(http.StatusNotFound, errorResponse)
			}
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to fetch medical record"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		successResponse := map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Medical record fetched successfully",
			"data":    medicalRecord,
		}
		return c.JSON(http.StatusOK, successResponse)
	}
}

func EditMedicalRecordByID(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		doctor := c.Get("doctor").(*models.Doctor)

		medicalRecordID := c.Param("id")

		var existingMedicalRecord models.MedicalRecords
		result := db.First(&existingMedicalRecord, medicalRecordID)
		if result.Error != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Medical record not found",
			}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		if existingMedicalRecord.DoctorID != doctor.ID {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "You are not authorized to edit this medical record",
			}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		var updatedMedicalRecord models.MedicalRecords
		if err := c.Bind(&updatedMedicalRecord); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		if updatedMedicalRecord.PatientName != "" {
			if len(updatedMedicalRecord.PatientName) < 1 || len(updatedMedicalRecord.PatientName) > 100 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "PatientName must be between 1 and 100 characters long",
				})
			}
			if !helper.ValidateLettersAndSpaces(updatedMedicalRecord.PatientName) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Invalid patient name format. Only letters and spaces are allowed",
				})
			}
			existingMedicalRecord.PatientName = updatedMedicalRecord.PatientName
		}

		if updatedMedicalRecord.BirthDate != "" {
			if _, err := time.Parse("2006-01-02", updatedMedicalRecord.BirthDate); err != nil {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "BirthDate must be in the format yyyy-mm-dd",
				})
			}
			existingMedicalRecord.BirthDate = updatedMedicalRecord.BirthDate
		}

		if updatedMedicalRecord.Email != "" {
			if !helper.ValidateEmailFormat(updatedMedicalRecord.Email) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Email must be a valid email format",
				})
			}
			existingMedicalRecord.Email = updatedMedicalRecord.Email
		}

		if updatedMedicalRecord.PhoneNumber != "" {
			if len(updatedMedicalRecord.PhoneNumber) > 13 || !helper.ValidatePhoneNumber(updatedMedicalRecord.PhoneNumber) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "PhoneNumber must contain only digits and be at most 13 characters long",
				})
			}
			existingMedicalRecord.PhoneNumber = updatedMedicalRecord.PhoneNumber
		}

		if updatedMedicalRecord.Diagnosis != "" {
			if len(updatedMedicalRecord.Diagnosis) < 5 || len(updatedMedicalRecord.Diagnosis) > 3000 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Diagnosis must be between 5 and 3000 characters long",
				})
			}
			existingMedicalRecord.Diagnosis = updatedMedicalRecord.Diagnosis
		}

		if updatedMedicalRecord.Prescription != "" {
			if len(updatedMedicalRecord.Prescription) < 5 || len(updatedMedicalRecord.Prescription) > 3000 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Prescription must be between 5 and 3000 characters long",
				})
			}
			existingMedicalRecord.Prescription = updatedMedicalRecord.Prescription
		}

		if updatedMedicalRecord.CareSuggestion != "" {
			if len(updatedMedicalRecord.CareSuggestion) < 5 || len(updatedMedicalRecord.CareSuggestion) > 3000 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "CareSuggestion must be between 5 and 3000 characters long",
				})
			}
			existingMedicalRecord.CareSuggestion = updatedMedicalRecord.CareSuggestion
		}

		db.Save(&existingMedicalRecord)

		medicalRecordResponse := map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Medical record updated successfully",
			"data": map[string]interface{}{
				"id":              existingMedicalRecord.ID,
				"patient_name":    existingMedicalRecord.PatientName,
				"birth_date":      existingMedicalRecord.BirthDate,
				"email":           existingMedicalRecord.Email,
				"phone_number":    existingMedicalRecord.PhoneNumber,
				"diagnosis":       existingMedicalRecord.Diagnosis,
				"prescription":    existingMedicalRecord.Prescription,
				"care_suggestion": existingMedicalRecord.CareSuggestion,
				"created_at":      existingMedicalRecord.CreatedAt,
				"updated_at":      existingMedicalRecord.UpdatedAt,
			},
		}

		return c.JSON(http.StatusOK, medicalRecordResponse)
	}
}

func DeleteMedicalRecordByID(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		doctor := c.Get("doctor").(*models.Doctor)

		medicalRecordID := c.Param("id")

		var existingMedicalRecord models.MedicalRecords
		result := db.First(&existingMedicalRecord, medicalRecordID)
		if result.Error != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Medical record not found",
			}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		if existingMedicalRecord.DoctorID != doctor.ID {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "You are not authorized to delete this medical record",
			}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		if err := db.Delete(&existingMedicalRecord).Error; err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete medical record",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		successResponse := map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Medical record deleted successfully",
		}
		return c.JSON(http.StatusOK, successResponse)
	}
}

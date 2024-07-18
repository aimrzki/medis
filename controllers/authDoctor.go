package controllers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"medis/helper"
	"medis/middleware"
	"medis/models"
	"net/http"
	"regexp"
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
		var doctor models.Doctor
		if err := c.Bind(&doctor); err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(doctor.FirstName) < 3 || len(doctor.LastName) < 3 {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "First name and last name must be at least 3 characters"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(doctor.Username) < 5 {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Username must be at least 5 characters"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if len(doctor.Password) < 8 || !helper.IsValidPassword(doctor.Password) {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Password must be at least 8 characters and contain a combination of letters and numbers"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		match, _ := regexp.MatchString(emailPattern, doctor.Email)
		if !match {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid email format"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		var existingDoctor models.Doctor
		result := db.Where("username = ?", doctor.Username).First(&existingDoctor)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusConflict, Message: "Username already exists"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check username"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		result = db.Where("email = ?", doctor.Email).First(&existingDoctor)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusConflict, Message: "Email already exists"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to hash password"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		uniqueToken := helper.GenerateUniqueToken()
		doctor.VerificationToken = uniqueToken

		doctor.Password = string(hashedPassword)
		doctor.Fullname = doctor.FirstName + " " + doctor.LastName
		db.Create(&doctor)

		doctor.Password = ""

		tokenString, err := middleware.GenerateToken(doctor.Username, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		if err := helper.SendWelcomeEmail(doctor.Email, doctor.FirstName+" "+doctor.LastName, uniqueToken); err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to send welcome email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
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
		var doctor models.Doctor
		if err := c.Bind(&doctor); err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}
		if doctor.Username == "" {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Username is required"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}
		if doctor.Password == "" {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Password is required"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		var existingDoctor models.Doctor
		result := db.Where("username = ?", doctor.Username).First(&existingDoctor)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorResponse := helper.ErrorResponse{Code: http.StatusUnauthorized, Message: "Invalid username"}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			} else {
				errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check username"}
				return c.JSON(http.StatusInternalServerError, errorResponse)
			}
		}

		err := bcrypt.CompareHashAndPassword([]byte(existingDoctor.Password), []byte(doctor.Password))
		if err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusUnauthorized, Message: "Invalid password"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		if !existingDoctor.IsVerified {
			errorResponse := helper.ErrorResponse{Code: http.StatusUnauthorized, Message: "Account not verified. Please verify your account before logging in."}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		tokenString, err := middleware.GenerateToken(existingDoctor.Username, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		go func(email, username string) {
			if err := helper.SendLoginNotification(email, username); err != nil {
				fmt.Println("Failed to send notification email:", err)
			}
		}(existingDoctor.Email, existingDoctor.FirstName+" "+existingDoctor.LastName)

		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Doctor login successful", "token": tokenString, "id": existingDoctor.ID})
	}
}

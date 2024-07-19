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

/*
Function VerifyEmail berfungsi untuk memverivikasi token yang dikirimkan ke
dalam email dokter yang mendaftar, jika dokter telah mengklik tombol tersebut
maka status dokter akan berubah menjadi verified dan telah diizinkan untuk login
*/

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

/*
Function RegisterDoctor digunakan untuk melakukan pendaftaran akun dokter di platform
pada endpoint ini ada beberapa data yang harus dimasukan seperti
(firstname, lastname, email, username, dan password), saat dokter melakukan pendaftaran
password akan dihash menggunakan library bacrypt, lalu setelah semya data valid maka data akan disimpan
ke dalam table doctors
*/
func RegisterDoctor(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var doctor models.Doctor
		if err := c.Bind(&doctor); err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan first name minimal terdiri 1 huruf dan maksimal 100 huruf
		if len(doctor.FirstName) < 1 || len(doctor.FirstName) > 100 || !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(doctor.FirstName) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "First Name must be between 1 and 100 characters and contain only letters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan last name boleh kosong dan bila diisi maksimal 100 huruf
		if len(doctor.LastName) > 0 {
			if len(doctor.LastName) > 100 || !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(doctor.LastName) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Last name max 100 characters and contain only letters",
				})
			}
		}

		// Menyimpan nilai lastname kedalam database
		doctor.LastName = doctor.LastName

		// Logika untuk memastikan last username boleh kosong dan bila diisi maksimal 100 huruf
		if len(doctor.Username) < 5 || len(doctor.Username) > 100 {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Username must be at least 5 characters and max 100 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan last password minimal 8 karakter maksimal 100 karakter
		if len(doctor.Password) < 8 || len(doctor.Password) > 100 || !helper.IsValidPassword(doctor.Password) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Password must be at least 8 characters max 100 characters and contain a combination of letters and numbers",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan email dimasukan dengan format yang benar
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		match, _ := regexp.MatchString(emailPattern, doctor.Email)
		if !match {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid email format",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan bahwa nomor telphon hanya bisa angka dan minimal 10 angka maksimal 13 angka sesuai format nomor di indonesia
		contactNumberRegex := regexp.MustCompile(`^\d{10,13}$`)
		if !contactNumberRegex.MatchString(doctor.ContactNumber) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Contact number must be between 10 and 14 digits and contain only numbers",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk melalukan pengecekan apakah username yang diinputkan tersedia belum di pakai oleh dokter lain
		var existingDoctor models.Doctor
		result := db.Where("username = ?", doctor.Username).First(&existingDoctor)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Username already exists",
			}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to check username",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Logika untuk melalukan pengecekan apakah email yang diinputkan tersedia belum di pakai oleh dokter lain
		result = db.Where("email = ?", doctor.Email).First(&existingDoctor)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Email already exists",
			}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to check email",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Kode untuk melakukan hashing password menggunakan library bacrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to hash password",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Kode untuk menghasilakn token yang akan dikirim ke email untuk verifikasi akun dokter
		uniqueToken := helper.GenerateUniqueToken()
		doctor.VerificationToken = uniqueToken

		doctor.Password = string(hashedPassword)
		doctor.Fullname = doctor.FirstName + " " + doctor.LastName
		db.Create(&doctor)

		doctor.Password = ""

		tokenString, err := middleware.GenerateToken(doctor.Username, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to generate token",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		if err := helper.SendWelcomeEmail(doctor.Email, doctor.FirstName+" "+doctor.LastName, uniqueToken); err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to send welcome email",
			}
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

/*
Function SignInDoctor digunakan untuk endpoint dokter melakukan signin ke dalam sistem
dokter diharuskan untuk mengisi username, password, dan memastikan telah melakukan verifikasi
pada link yang dikirim ke email untuk dapat login kedalam sistem
*/

func SignInDoctor(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var doctor models.Doctor
		if err := c.Bind(&doctor); err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan username yang diinputkan tidak boleh kosong
		if doctor.Username == "" {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Username is required",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk memastikan password yang diinputkan tidak boleh kosong
		if doctor.Password == "" {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Password is required",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Logika untuk mengecek username yang diinputkan apakah tersedia di dalam database
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

		// Logika untuk mengecek password yang diinputkan apakah sesuai dengan password di database yang telah di hash
		err := bcrypt.CompareHashAndPassword([]byte(existingDoctor.Password), []byte(doctor.Password))
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid password",
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Logika untuk memastikan dokter yang bisa masuk ke dalam sistem hanya yang sudah verifikasi link yang dikirim ke email saat pendaftaran
		if !existingDoctor.IsVerified {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Account not verified. Please verify your account before logging in.",
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Logika untuk menghasilkan token authorization (bearer  token) untuk mengakses semua fitur dokter
		tokenString, err := middleware.GenerateToken(existingDoctor.Username, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to generate token",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Kode untuk mengirimkan notifikasi peringkatan login kepada email dokter
		go func(email, username string) {
			if err := helper.SendLoginNotification(email, username); err != nil {
				fmt.Println("Failed to send notification email:", err)
			}
		}(existingDoctor.Email, existingDoctor.FirstName+" "+existingDoctor.LastName)

		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "Doctor login successful", "token": tokenString, "id": existingDoctor.ID})
	}
}

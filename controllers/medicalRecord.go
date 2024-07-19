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

/*
Function AddMedicalRecord berfungsi untuk endpoint agar dokter dapat menambahkan riwayar medical record
untuk suatu pasien dengan memasukan nama pasien, tanggal lahir, email, nomor telphone, diagnosis , resep obat, dan saran perawatan
*/
func AddMedicalRecord(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		/*
			Function yang berguna untuk mengecek token authorization yang dimasukan,
			apakah sah dokter atau tidak
		*/
		doctor, err := helper.VerifyDoctorToken(db, c, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		/*
			Kode untuk menginputkan data yang diperlukan ke dalam
			tabel medical records
		*/
		var medicalRecord models.MedicalRecords
		if err := c.Bind(&medicalRecord); err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika untuk memastikan bahwa patient name minimal harus memiliki
			1 huruf dan maksimal 100 huruf
		*/
		if len(medicalRecord.PatientName) < 1 || len(medicalRecord.PatientName) > 100 || !helper.ValidateLettersAndSpaces(medicalRecord.PatientName) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Patient name must be between 1 and 100 characters and contain only letters and spaces",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika yang digunakan untuk mengecek inputan tanggal lahir
			format harus dalam bentuk (YYYY-MM-DD)
			contoh : 2001-11-11
		*/
		if !helper.ValidateDateFormat(medicalRecord.BirthDate) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Birth date must be in the format yyyy-mm-dd",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika yang digunakan untuk mengecek inputan email yang dimasukan
			format harus sesuai dengan inputan email yang seharusnya
			contoh : user@gmail.com
		*/
		if !helper.ValidateEmailFormat(medicalRecord.Email) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid email format",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika yang digunakan untuk mengecek inputan nomor telphone yang dimasukan
			inputan harus hanya berupa angka
			dengan minimal 10 digit dan maksimal 13 digit
		*/
		if !helper.ValidatePhoneNumber(medicalRecord.PhoneNumber) {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Phone number must contain only digits and be at most 13 characters long",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika yang digunakan untuk mengecek inputan bahwa diagnosis harus
			minimal 5 karakter dan maksimal 3000 karakter
		*/
		if len(medicalRecord.Diagnosis) < 1 || len(medicalRecord.Diagnosis) > 3000 {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest,
				Message: "Diagnosis must be between 1 and 3000 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika yang digunakan untuk mengecek inputan bahwa resep obat harus
			minimal 1 karakter dan maksimal 3000 karakter
		*/
		if len(medicalRecord.Prescription) < 1 || len(medicalRecord.Prescription) > 3000 {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest,
				Message: "Prescription must be between 1 and 3000 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			Logika yang digunakan untuk mengecek inputan bahwa saran perawatan harus
			minimal 5 karakter dan maksimal 3000 karakter
		*/
		if len(medicalRecord.CareSuggestion) < 1 || len(medicalRecord.CareSuggestion) > 3000 {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Care suggestion must be between 1 and 3000 characters",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		/*
			birthDate, err := time.Parse("2006-01-02", medicalRecord.BirthDate)
			if err != nil {
				errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid BirthDayDate format"}
				return c.JSON(http.StatusBadRequest, errorResponse)
			}
			medicalRecord.BirthDate = birthDate.Format("2006-01-02")
		*/

		// Variable untuk menyimpan dokter id siapa yang membuat data ini
		medicalRecord.DoctorID = doctor.ID

		/*
			Logika untuk melakukan penyimpanan data ke dalam table medical record
		*/
		if err := db.Create(&medicalRecord).Error; err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create medical record",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		/*
			Kode untuk mengirimkan pesan notifikasi ke email pasien yang diinputkan
			sekaligus mengirimkan pdf hasil medical record
		*/
		if err := helper.SendMedicalRecordNotification(medicalRecord.Email, medicalRecord.PatientName, medicalRecord.Diagnosis, medicalRecord.Prescription, medicalRecord.CareSuggestion); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to send medical record notification",
			})
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

/*
Function GetMedicalRecordsByDoctor digunakan oleh dokter untuk melihat riwayat medical record
untuk pasian yang telah dia tambahkan, dokter hanya dapat melihat data medical record yang dirinya tambahkan
tidak bisa yang di tambahkan oleh dokter lain
*/

func GetMedicalRecordsByDoctor(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Logika untuk memverifikasi token yang dimasukan apakah valid milik dokter
		doctor, err := helper.VerifyDoctorToken(db, c, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Logika untuk pagination dengan membaca query params page dan limit
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}

		// Menampilkan data medical record dari table
		var medicalRecords []models.MedicalRecords
		offset := (page - 1) * limit
		if err := db.Where("doctor_id = ?", doctor.ID).Offset(offset).Limit(limit).Find(&medicalRecords).Error; err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to fetch medical records",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Logika untuk menghitung total record yang ada pada medical record
		var totalRecords int64
		db.Model(&models.MedicalRecords{}).Where("doctor_id = ?", doctor.ID).Count(&totalRecords)

		// Response saat berhasil
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

/*
Function GetMedicalRecordByID digunakan untuk dokter dapat menampilkan detail data medical record
milik pasien berdasarkan medical record ID nya
*/
func GetMedicalRecordByID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Logika untuk memverifikasi token yang dimasukan apakah valid milik dokter
		doctor, err := helper.VerifyDoctorToken(db, c, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Pengecekan inputan id yang dikirim pada param
		recordID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid record ID",
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Pengecekan ID apakah tersedia di databases dan aksesnya
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

		//Response saat berhasil
		successResponse := map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Medical record fetched successfully",
			"data":    medicalRecord,
		}
		return c.JSON(http.StatusOK, successResponse)
	}
}

/*
Function EditMedicalRecordByID digunakan untuk dokter dapat mengubah data medical record berdasarkan
ID nya, dokter dapat mengubah hanya salah satu bagiannya saja atau bahkan bisa semuanya.
Tinggal sesuaikan saja key dan valuenya yang ingin diedit.
*/

func EditMedicalRecordByID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Logika untuk memverifikasi token yang dimasukan apakah valid milik dokter
		doctor, err := helper.VerifyDoctorToken(db, c, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Code untuk menerima id dari parameter
		medicalRecordID := c.Param("id")

		// Logika pengecekan apakah medical record tersedia di databases
		var existingMedicalRecord models.MedicalRecords
		result := db.First(&existingMedicalRecord, medicalRecordID)
		if result.Error != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Medical record not found",
			}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Logika pengecekan apakah dokter id pada medical record sesuai dengan dokter id yang login
		if existingMedicalRecord.DoctorID != doctor.ID {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "You are not authorized to edit this medical record",
			}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Logika untuk pengecekan inputan yang di masukan apakah sesuai
		var updatedMedicalRecord models.MedicalRecords
		if err := c.Bind(&updatedMedicalRecord); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		/*
			Logika perubahan nama pasien, jika nama pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usdah dikirim, namun jika ingin diubah pastikan nama pasien terdiri dari
			minimal 1 huruf dan maksimal 1000 huruf
		*/
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

		/*
			Logika perubahan tanggal lahir pasien, jika tanggal lahir pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usdah dikirim, namun jika ingin diubah pastikan tanggal lahir pasien berformat YYYY-MM-DD
		*/
		if updatedMedicalRecord.BirthDate != "" {
			if _, err := time.Parse("2006-01-02", updatedMedicalRecord.BirthDate); err != nil {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "BirthDate must be in the format yyyy-mm-dd",
				})
			}
			existingMedicalRecord.BirthDate = updatedMedicalRecord.BirthDate
		}

		/*
			Logika perubahan email pasien, jika email pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usdah dikirim, namun jika ingin diubah pastikan email memiliki format yang sesuai
		*/
		if updatedMedicalRecord.Email != "" {
			if !helper.ValidateEmailFormat(updatedMedicalRecord.Email) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Email must be a valid email format",
				})
			}
			existingMedicalRecord.Email = updatedMedicalRecord.Email
		}

		/*
			Logika perubahan nomor telphone pasien, jika nomor telphone pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usdah dikirim, namun jika ingin diubah pastikan nomor telphone terdiri dari hanya angka minimal 10 maksimal 13 digit
		*/
		if updatedMedicalRecord.PhoneNumber != "" {
			if len(updatedMedicalRecord.PhoneNumber) > 13 || !helper.ValidatePhoneNumber(updatedMedicalRecord.PhoneNumber) {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "PhoneNumber must contain only digits and be at most 13 characters long",
				})
			}
			existingMedicalRecord.PhoneNumber = updatedMedicalRecord.PhoneNumber
		}

		/*
			Logika perubahan diagnosis pasien, jika diagnosise pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usah dikirim, namun jika ingin diubah pastikan diagnosise terdiri dari minimal 5 karakter maksimal 3000 karakter
		*/
		if updatedMedicalRecord.Diagnosis != "" {
			if len(updatedMedicalRecord.Diagnosis) < 5 || len(updatedMedicalRecord.Diagnosis) > 3000 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Diagnosis must be between 5 and 3000 characters long",
				})
			}
			existingMedicalRecord.Diagnosis = updatedMedicalRecord.Diagnosis
		}

		/*
			Logika perubahan resep obat pasien, jika resep obat pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usah dikirim, namun jika ingin diubah pastikan resep obat terdiri dari minimal 5 karakter maksimal 3000 karakter
		*/
		if updatedMedicalRecord.Prescription != "" {
			if len(updatedMedicalRecord.Prescription) < 5 || len(updatedMedicalRecord.Prescription) > 3000 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Prescription must be between 5 and 3000 characters long",
				})
			}
			existingMedicalRecord.Prescription = updatedMedicalRecord.Prescription
		}

		/*
			Logika perubahan saran perawatan pasien, jika saran perawatan pasien tidak mau di ubah cukup kosongkan saja
			atau tidak usah dikirim, namun jika ingin diubah pastikan saran perawatan terdiri dari minimal 5 karakter maksimal 3000 karakter
		*/
		if updatedMedicalRecord.CareSuggestion != "" {
			if len(updatedMedicalRecord.CareSuggestion) < 5 || len(updatedMedicalRecord.CareSuggestion) > 3000 {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "CareSuggestion must be between 5 and 3000 characters long",
				})
			}
			existingMedicalRecord.CareSuggestion = updatedMedicalRecord.CareSuggestion
		}

		// kode untuk melakukan perubahan ke databases sesuai dengan data yang dikirim
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

func DeleteMedicalRecordByID(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Logika untuk memverifikasi token yang dimasukan apakah valid milik dokter
		doctor, err := helper.VerifyDoctorToken(db, c, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Kode untuk menerima id yang dikirimkan dari param
		medicalRecordID := c.Param("id")

		// Logika pengecekan medical record pada databases
		var existingMedicalRecord models.MedicalRecords
		result := db.First(&existingMedicalRecord, medicalRecordID)
		if result.Error != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Medical record not found",
			}
			return c.JSON(http.StatusNotFound, errorResponse)
		}

		// Logika pengecekan apakah dokter yang ingin menghapus sesuai dengan dokter id pada medical record
		if existingMedicalRecord.DoctorID != doctor.ID {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "You are not authorized to delete this medical record",
			}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		// Logika untuk melakukan proses penghapusan medical record dari databases
		if err := db.Delete(&existingMedicalRecord).Error; err != nil {
			errorResponse := helper.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete medical record",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// response jika berhasil
		successResponse := map[string]interface{}{
			"code":    http.StatusOK,
			"error":   false,
			"message": "Medical record deleted successfully",
		}
		return c.JSON(http.StatusOK, successResponse)
	}
}

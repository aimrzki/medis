package controllers

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"medis/helper"
	"net/http"
	"os"
	"strconv"
)

// Fungsi untuk mendapatkan token autentikasi
func GetAuthToken(c echo.Context) error {
	// Bind data dari request ke dalam struct AuthRequest
	var authReq helper.AuthRequest
	if err := c.Bind(&authReq); err != nil {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to bind request body: " + err.Error(),
		}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	// Mengambil URL dari variabel lingkungan
	authURL := os.Getenv("AUTH_URL")
	if authURL == "" {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "AUTH_URL is not set in the environment",
		}
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	// Mengirim request ke API untuk mendapatkan token autentikasi
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     authReq.ClientID,
			"client_secret": authReq.ClientSecret,
			"grant_type":    authReq.GrantType,
		}).
		SetQueryParam("grant_type", authReq.GrantType).
		Post(authURL)

	if err != nil {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get access token: " + err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	var authResponse helper.AuthResponse
	if err := json.Unmarshal(resp.Body(), &authResponse); err != nil {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse access token response: " + err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	return c.JSON(http.StatusOK, authResponse)
}

// Fungsi untuk mendapatkan daftar obat
func GetMedicineList(c echo.Context) error {

	// Kode untuk memverifikasi token yang diinputkan
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Token is required",
		}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	// Page limit pada query param untuk pagination
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	// Jika page dan limit tidak di masukan defaultnya page 1 dan limit 10
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Mengambil URL dari variabel lingkungan
	medicineURL := os.Getenv("MEDICINE_URL")
	if medicineURL == "" {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "MEDICINE_URL is not set in the environment",
		}
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	// Mengirimkan rest API untuk melihat daftar obat
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", authHeader).
		SetQueryParams(map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		}).
		Get(medicineURL)

	if err != nil {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get medicine list: " + err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	var medicineResponse helper.MedicineResponse
	if err := json.Unmarshal(resp.Body(), &medicineResponse); err != nil {
		errorResponse := helper.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse medicine list response: " + err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errorResponse)
	}

	return c.JSON(http.StatusOK, medicineResponse)
}

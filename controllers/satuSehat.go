package controllers

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strconv"
)

type AuthResponse struct {
	RefreshTokenExpiresIn string   `json:"refresh_token_expires_in"`
	ApiProductList        string   `json:"api_product_list"`
	ApiProductListJson    []string `json:"api_product_list_json"`
	OrganizationName      string   `json:"organization_name"`
	DeveloperEmail        string   `json:"developer.email"`
	TokenType             string   `json:"token_type"`
	IssuedAt              string   `json:"issued_at"`
	ClientID              string   `json:"client_id"`
	AccessToken           string   `json:"access_token"`
	ApplicationName       string   `json:"application_name"`
	Scope                 string   `json:"scope"`
	ExpiresIn             string   `json:"expires_in"`
	RefreshCount          string   `json:"refresh_count"`
	Status                string   `json:"status"`
}

type MedicineResponse struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Items struct {
		Data []Medicine `json:"data"`
	} `json:"items"`
}

type Medicine struct {
	KfaCode             *string    `json:"kfa_code"`
	ProductTemplateName *string    `json:"product_template_name"`
	DocumentRef         string     `json:"document_ref"`
	Active              bool       `json:"active"`
	RegionName          string     `json:"region_name"`
	RegionCode          string     `json:"region_code"`
	StartDate           string     `json:"start_date"`
	EndDate             *string    `json:"end_date"`
	PriceUnit           int        `json:"price_unit"`
	UomName             *string    `json:"uom_name"`
	UpdatedAt           string     `json:"updated_at"`
	UomPack             []string   `json:"uom_pack"`
	Province            []Province `json:"province"`
}

type Province struct {
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

func GetAuthToken(c echo.Context) error {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"grant_type":    "client_credentials",
		}).
		Post("https://api-satusehat-stg.dto.kemkes.go.id/oauth2/v1/accesstoken?grant_type=client_credentials")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get access token: " + err.Error()})
	}

	var authResponse AuthResponse
	if err := json.Unmarshal(resp.Body(), &authResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse access token response: " + err.Error()})
	}

	return c.JSON(http.StatusOK, authResponse)
}

/*
func GetAuthToken(c echo.Context) error {
	type AuthRequest struct {
		ClientID     string `form:"client_id"`
		ClientSecret string `form:"client_secret"`
		GrantType    string `form:"grant_type"`
	}

	var authReq AuthRequest
	if err := c.Bind(&authReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to bind request body: " + err.Error()})
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"client_id":     authReq.ClientID,
			"client_secret": authReq.ClientSecret,
			"grant_type":    authReq.GrantType,
		}).
		SetQueryParam("grant_type", authReq.GrantType).
		Post("https://api-satusehat-stg.dto.kemkes.go.id/oauth2/v1/accesstoken")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get access token: " + err.Error()})
	}

	var authResponse AuthResponse
	if err := json.Unmarshal(resp.Body(), &authResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse access token response: " + err.Error()})
	}

	return c.JSON(http.StatusOK, authResponse)
}
*/

func GetMedicineList(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}
	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", authHeader).
		SetQueryParams(map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		}).
		Get("https://api-satusehat-stg.kemkes.go.id/kfa/farmalkes-price-jkn")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get medicine list: " + err.Error()})
	}

	var medicineResponse MedicineResponse
	if err := json.Unmarshal(resp.Body(), &medicineResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse medicine list response: " + err.Error()})
	}

	return c.JSON(http.StatusOK, medicineResponse)
}

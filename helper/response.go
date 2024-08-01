package helper

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Struktur untuk request autentikasi
type AuthRequest struct {
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	GrantType    string `form:"grant_type"`
}

// Struktur untuk menampung response dari API Auth
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

// Struktur untuk menampung response dari API daftar obat
type MedicineResponse struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Items struct {
		Data []Medicine `json:"data"`
	} `json:"items"`
}

// Struktur untuk data obat
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

// Struktur untuk data provinsi
type Province struct {
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

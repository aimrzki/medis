package models

type Doctor struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Fullname          string `json:"fullname"`
	ContactNumber     string `json:"contact_number"`
	Gender            string `json:"gender"`
	Email             string `json:"email"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	IsVerified        bool   `gorm:"default:false" json:"is_verified"`
	VerificationToken string `json:"verification_token"`
}

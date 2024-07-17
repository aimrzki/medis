package models

import "time"

type MedicalRecords struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	PatientName    string     `json:"patient_name"`
	BirthDate      string     `json:"birth_date"`
	Email          string     `json:"email"`
	PhoneNumber    string     `json:"phone_number"`
	Diagnosis      string     `json:"diagnosis"`
	Prescription   string     `json:"prescription"`
	CareSuggestion string     `json:"care_suggestion"`
	DoctorID       uint       `json:"doctor_id"` // Foreign key to Doctor
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      time.Time
}

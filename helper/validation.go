package helper

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"time"
)

func IsValidPassword(password string) bool {
	// Password must be at least 8 characters and contain a combination of letters and numbers
	return len(password) >= 8 && containsLetterAndNumber(password)
}

func containsLetterAndNumber(s string) bool {
	hasLetter := false
	hasNumber := false
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			return true
		}
	}
	return false
}

func GenerateUniqueToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func ValidateLettersAndSpaces(input string) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return re.MatchString(input)
}

func ValidateDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func ValidateEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidatePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\d{10,13}$`)
	return re.MatchString(phone)
}

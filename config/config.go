package config

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"medis/models"
	"os"
	"strconv"
)

// Object database yang harus diisikan berupa (DB HOST, DB PORT, DB USERNAME, DB PASSWORD, dan DB BAME)
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

/*
Function untuk mengkonfigurasi database, semua nilai yang dibutuhkan untuk database seperti
(DB HOST, DB PORT, DB USERNAME, DB PASSWORD, dan DB BAME) diambil langsung dari .env
pada CI/CD nilai tersebut diambil dari github secret
*/
func InitializeDatabase() (*gorm.DB, error) {
	godotenv.Load(".env")

	dbConfig := DatabaseConfig{
		Host: os.Getenv("DB_HOST"),
	}
	portStr := os.Getenv("DB_PORT")
	dbConfig.Port, _ = strconv.Atoi(portStr)
	dbConfig.Username = os.Getenv("DB_USERNAME")
	dbConfig.Password = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")

	dsn := "user=" + dbConfig.Username + " password=" + dbConfig.Password + " dbname=" + dbConfig.DBName + " host=" + dbConfig.Host + " port=" + strconv.Itoa(dbConfig.Port) + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Gorm yang digunakan untuk konfigurasi posgree
	if err != nil {
		return nil, err
	}

	/*
		Kode untuk migrasi model model object ke dalam basis data menggunakan GORM
	*/
	db.AutoMigrate(&models.Doctor{})
	db.AutoMigrate(&models.MedicalRecords{})

	return db, nil
}

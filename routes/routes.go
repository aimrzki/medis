package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"io/ioutil"
	"medis/controllers"
	"medis/middleware"
	"net/http"
)

func ServeHTML(c echo.Context) error {
	htmlData, err := ioutil.ReadFile("index.html")
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, string(htmlData))
}

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	e.Use(Logger())
	secretKey := []byte(middleware.GetSecretKeyFromEnv())
	e.GET("/", ServeHTML)

	e.POST("/api/doctor/signup", controllers.RegisterDoctor(db, secretKey))
	e.POST("/api/doctor/signin", controllers.SignInDoctor(db, secretKey))
	e.GET("/verify", controllers.VerifyEmail(db))

	// Satu Sehat
	e.POST("/api/satusehat/auth", controllers.GetAuthToken)
	e.GET("/api/satusehat/medicine", controllers.GetMedicineList)

	// Medical Record
	e.POST("/api/doctor/medical-record", controllers.AddMedicalRecord(db, secretKey))
	e.GET("/api/doctor/medical-record", controllers.GetMedicalRecordsByDoctor(db, secretKey))
	e.GET("/api/doctor/medical-record/:id", controllers.GetMedicalRecordByID(db, secretKey))
	e.PUT("/api/doctor/medical-record/:id", controllers.EditMedicalRecordByID(db, secretKey))
	e.DELETE("/api/doctor/medical-record/:id", controllers.DeleteMedicalRecordByID(db, secretKey))

}

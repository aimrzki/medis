package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"io/ioutil"
	"medis/auth"
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
	secretKey := []byte(auth.GetSecretKeyFromEnv())
	e.GET("/", ServeHTML)

	e.POST("/api/doctor/signup", middleware.ValidateDoctorRegistration(middleware.CheckDoctorUniqueness(db)(controllers.RegisterDoctor(db, secretKey))))
	e.POST("/api/doctor/signin", middleware.ValidateDoctorSignIn(db)(controllers.SignInDoctor(db, secretKey)))
	e.GET("/verify", controllers.VerifyEmail(db))

	// Satu Sehat
	e.POST("/api/satusehat/auth", controllers.GetAuthToken)
	e.GET("/api/satusehat/medicine", controllers.GetMedicineList)

	// Medical Record
	e.POST("/api/doctor/medical-record",
		middleware.VerifyDoctorTokenMiddleware(db, secretKey)(
			middleware.ValidateMedicalRecord(
				controllers.AddMedicalRecord(db),
			),
		),
	)
	e.GET("/api/doctor/medical-record",
		middleware.VerifyDoctorTokenMiddleware(db, secretKey)(
			controllers.GetMedicalRecordsByDoctor(db),
		),
	)

	e.GET("/api/doctor/medical-record/:id",
		middleware.VerifyDoctorTokenMiddleware(db, secretKey)(
			controllers.GetMedicalRecordByID(db),
		),
	)

	e.PUT("/api/doctor/medical-record/:id",
		middleware.VerifyDoctorTokenMiddleware(db, secretKey)(
			controllers.EditMedicalRecordByID(db),
		),
	)

	e.DELETE("/api/doctor/medical-record/:id",
		middleware.VerifyDoctorTokenMiddleware(db, secretKey)(
			controllers.DeleteMedicalRecordByID(db),
		),
	)
	
}

package config

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"medis/routes"
)

/*
Function yang digunakan untuk mensetup middleware yang digunakan
middleware logger -> digunakan untuk mencatat logging aktivitas yang terjadi di dalam aplikasi
middleware cors -> digunakan untuk memungkinkan aplikasi untuk berkomunikasi dengan API atau layanan web lain yang mungkin berada di domain yang berbeda
middleware RemoveTrailingSlash -> digunakan untuk menghapus otomatis tanda garis miring (/) di akhir URL yang diminta. Misalnya, jika ada permintaan ke /about/, middleware ini akan secara otomatis mengarahkannya ke /about.
*/

func SetupRouter() *echo.Echo {
	db, err := InitializeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Pre(middleware.RemoveTrailingSlash())
	routes.SetupRoutes(router, db)
	return router
}

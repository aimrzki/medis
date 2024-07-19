package main

import (
	"log"
	"medis/config"
)

/*
Function main untuk memulai server pada port 8080
*/
func main() {
	router := config.SetupRouter()
	err := router.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

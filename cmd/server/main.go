package main

import (
	"DropFiles/internal/app"
	"log"
)

// @title DropFiles API
// @version 1.0
// @description File sharing service like Dropmefiles.com. Upload files without auth, get download link.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

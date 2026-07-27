package main

import (
	"bms2/internal/app"
	"log"
)

func main() {

	app := app.New()

	log.Fatal(app.Run())
}

package main

import (
	"book-api/internal/app"
	"log"
)

func main() {

	app := app.New()

	log.Fatal(app.Run())
}

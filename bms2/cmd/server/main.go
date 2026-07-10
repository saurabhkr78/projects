package main

import (
	"log"

	"book-api/pkg/config"
	"book-api/pkg/database"
)

func main() {
	//laod the config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//connect to the database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Application started successfully.")

}

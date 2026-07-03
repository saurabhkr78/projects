package main

import (
	""
	"BMS/pkg/configs"
	"BMS/pkg/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//get the new router
	r := mux.NewRouter()
	//register the routes
	routes.RegisterBookRoutes(r)
	//connect to the database
	configs.Connect()
	//start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"fmt"
	"net/http"
	"redis-otp-service/client"
	"redis-otp-service/handler"

	"github.com/gorilla/mux"
)

func main() {
	//get a new https router
	router := mux.NewRouter()
	//initalize the redis client
	redisClient := client.NewClient()
	//Initialize the handlerfunc
	h := handler.NewHandler(redisClient)

	//register the handler functions for the routes
	router.HandleFunc("/request-otp", h.Requestotp).Methods("POST")
	router.HandleFunc("/verify-otp", h.Verifyotp).Methods("POST")

	//start the http server
	fmt.Println("starting the server @8080")
	http.ListenAndServe(":8080", router)

}

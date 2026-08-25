package main

import (
	"fmt"
	"net/http"
	"redis/internal/handler"
	client "redis/internal/redis"

	"github.com/gorilla/mux"
)

func main() {
	//initalize the redis client of redis go client from the client package
	redisClient := client.NewClient()

	// create a new router
	r := mux.NewRouter()

	//initalize the handler with the redis client
	h := handler.NewHandler(redisClient)

	// register the handler functions for the routes
	r.HandleFunc("/increment", h.Increment).Methods("POST")
	r.HandleFunc("/decrement", h.Decrement).Methods("POST")
	r.HandleFunc("/value", h.GetValue).Methods("GET")

	// start the http server
	fmt.Println("starting the server @8080")
	http.ListenAndServe(":8080", r)

}

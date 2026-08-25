package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"redis-rate-limiter/db"
	"redis-rate-limiter/handler"
	"redis-rate-limiter/middleware"
	"redis-rate-limiter/redis-client"
	"redis-rate-limiter/repository"
	"redis-rate-limiter/service"
)

func main() {
	//intialize the redis client
	redisClient := redisclient.NewClient()
	//initalize each layer of the application and inject the dependencies into each layer
	uRepo := repository.NewUserRepository(db.DB)

	uSvc := service.NewUserService(uRepo)

	//get new router
	r := mux.NewRouter()

	//register the handler with the router and inject the service into the handler

	uHandler := handler.NewHandler(uSvc)
	//initialize the rate limiter middleware and inject the redis client into the middleware
	rateLimit := middleware.NewLimiter(redisClient.RedisClient)
	//register the rate limiter middleware with the router
	r.Use(rateLimit.RateLimit)
	r.HandleFunc("/user/{id}", uHandler.Profile).Methods("GET")

	//start the server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

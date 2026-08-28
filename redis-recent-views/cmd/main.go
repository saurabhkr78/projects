package main

import (
	"context"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/saurabhkr78/redis-recent-views/db"
	"github.com/saurabhkr78/redis-recent-views/handler"
	"github.com/saurabhkr78/redis-recent-views/repository"
	"github.com/saurabhkr78/redis-recent-views/service"
)

func main() {
	// -------------------------
	// Fake DB
	// -------------------------
	products := db.Products

	// -------------------------
	// Redis
	// -------------------------
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis connection failed:", err)
	}

	// -------------------------
	// Repositories
	// -------------------------
	productRepo := repository.NewProductRepository(products)

	recentViewRepo := repository.NewRedisRecentViewRepository(
		redisClient,
	)

	// -------------------------
	// Service
	// -------------------------
	productService := service.NewProductService(
		productRepo,
		recentViewRepo,
	)

	// -------------------------
	// Handler
	// -------------------------
	h := handler.NewHandler(productService)

	// -------------------------
	// Router
	// -------------------------
	mux := http.NewServeMux()

	mux.HandleFunc("/products/", h.GetProduct)
	mux.HandleFunc("/users/recently-viewed", h.GetRecentViews)
	mux.HandleFunc("/users/recently-viewed/add", h.SetRecentView)

	// -------------------------
	// Server
	// -------------------------
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("server running on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

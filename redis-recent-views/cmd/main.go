package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
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
	productRepo := repository.NewFakeProductRepository(
		products,
	)

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
	// Fiber App
	// -------------------------
	app := fiber.New()

	// -------------------------
	// Routes
	// -------------------------

	// Get product.
	app.Get(
		"/products/:productID",
		h.GetProduct,
	)

	// Get recently viewed products.
	app.Get(
		"/users/:userID/recently-viewed",
		h.GetRecentViews,
	)

	// Add product to recently viewed.
	app.Post(
		"/users/:userID/recently-viewed/:productID",
		h.SetRecentView,
	)

	// -------------------------
	// Start server
	// -------------------------
	log.Println("server running on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}

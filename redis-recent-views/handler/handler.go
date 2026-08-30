package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saurabhkr78/redis-recent-views/service"
	"strconv"
)

type Handler struct {
	svc service.ProductService
}

func NewHandler(s service.ProductService) *Handler {
	return &Handler{
		svc: s,
	}
}

func (h *Handler) GetProduct(c *fiber.Ctx) error {
	// Get product ID from URL.
	productID := c.Params("productID")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product ID",
		})
	}

	// Call service.
	product, err := h.svc.GetProduct(
		c.Context(),
		productID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product",
		})
	}

	// Return product as JSON.
	return c.Status(fiber.StatusOK).JSON(product)
}

func (h *Handler) SetRecentView(c *fiber.Ctx) error {
	// Get user ID from URL.
	userID := c.Params("userID")

	// Get product ID from URL.
	productID := c.Params("productID")

	if userID == "" || productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID or product ID",
		})
	}

	// Call service.
	if err := h.svc.SetRecentView(
		c.Context(),
		userID,
		productID,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to set recent view",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetRecentViews(c *fiber.Ctx) error {
	// Get user ID from URL.
	userID := c.Params("userID")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	// Default limit.
	limit := 10

	// Get optional limit from query parameter.
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid limit",
			})
		}

		limit = parsedLimit
	}

	// Call service.
	productIDs, err := h.svc.GetRecentViews(
		c.Context(),
		userID,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get recent views",
		})
	}

	return c.Status(fiber.StatusOK).JSON(productIDs)
}

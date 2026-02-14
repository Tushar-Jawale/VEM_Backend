package handlers

import (
	"github.com/akashtripathi12/TBO_Backend/internal/models"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
	"github.com/gofiber/fiber/v2"
)

// GetHotelsByCity fetches the raw hotel list for a city (No rooms, just hotel info)
// GET /api/v1/hotels?city_id=DXB
func (r *Repository) GetHotelsByCity(c *fiber.Ctx) error {
	cityID := c.Query("city_id")

	// 1. Validate Input
	if cityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "city_id query parameter is required",
		})
	}

	var hotels []models.Hotel

	// 2. Query Database (Lightweight)
	// We SELECT * FROM hotels WHERE city_id = ?
	// We add a Limit(50) to prevent fetching 10,000 hotels at once
	result := store.DB.Where("city_id = ?", cityID).Limit(50).Find(&hotels)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch hotels",
		})
	}

	// 3. Return Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"count":  len(hotels),
		"data":   hotels,
	})
}

// GetHotel fetches a single hotel by its ID (hotel_code)
// GET /api/v1/hotels/:id
func (r *Repository) GetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	var hotel models.Hotel
	// Query by hotel_code (which is the primary key "ID" in struct, "hotel_code" in DB)
	result := store.DB.First(&hotel, "hotel_code = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Hotel not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   hotel,
	})
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/akashtripathi12/TBO_Backend/internal/models"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
	"github.com/akashtripathi12/TBO_Backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// GetHotelsByCity fetches the raw hotel list for a city with filters
// GET /api/v1/hotels?city_id=DXB&min_price=100&max_price=500&stars=4&rating=8.5&type=resort&hall_type=ballroom&hall_capacity=200
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
	// Generate cache key based on all query params to cache filtered results
	cacheKey := "hotels:city:" + cityID + ":" + c.OriginalURL()
	ctx := context.Background()

	// 2. Try to get from Redis
	if store.RDB != nil {
		cachedData, err := store.RDB.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedData), &hotels); err == nil {
				log.Printf("⚡ [REDIS] CACHE HIT: %s\n", cacheKey)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status": "success",
					"count":  len(hotels),
					"data":   hotels,
				})
			}
		}
	}

	// 3. Build Query
	query := store.DB.Model(&models.Hotel{}).Where("city_id = ?", cityID).
		Preload("Rooms").Preload("Banquets").Preload("Menus")

	// --- Room Filters ---
	joinRooms := false

	// 1. Free Cancellation
	if c.Query("free_cancellation") == "true" {
		joinRooms = true
		query = query.Where("room_offers.is_refundable = ?", true)
	}

	// 2. Room Configuration (Complex Filter)
	// ?room_config=[{"occupancy":2,"count":30},{"occupancy":3,"count":10}]
	if roomConfig := c.Query("room_config"); roomConfig != "" {
		type Req struct {
			Occupancy int `json:"occupancy"`
			Count     int `json:"count"`
		}
		var reqs []Req
		if err := json.Unmarshal([]byte(roomConfig), &reqs); err == nil {
			for _, r := range reqs {
				// Use a correlated subquery for each requirement
				// "Find hotels where EXISTS a room with cap>=X and count>=Y"
				subQuery := fmt.Sprintf("EXISTS (SELECT 1 FROM room_offers ro WHERE ro.hotel_id = hotels.hotel_code AND ro.max_capacity >= %d AND ro.count >= %d)", r.Occupancy, r.Count)
				query = query.Where(subQuery)
			}
		}
	}

	// 3. Basic Price/Guests Filters
	if minPrice := c.Query("min_price"); minPrice != "" {
		joinRooms = true
		query = query.Where("room_offers.total_fare >= ?", minPrice)
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		joinRooms = true
		query = query.Where("room_offers.total_fare <= ?", maxPrice)
	}
	if guestsPerRoom := c.Query("guests_per_room"); guestsPerRoom != "" {
		joinRooms = true
		query = query.Where("room_offers.max_capacity >= ?", guestsPerRoom)
	}
	if roomCount := c.Query("room_count"); roomCount != "" {
		joinRooms = true
		query = query.Where("room_offers.count >= ?", roomCount)
	}
	// Amenity Filter (JSONB array contains)
	// ?room_amenities=Bathtub,Balcony (comma separated)
	if roomAmenities := c.Query("room_amenities"); roomAmenities != "" {
		joinRooms = true
		amenitiesList := utils.ParseCSV(roomAmenities)
		for _, amenity := range amenitiesList {
			// Postgres JSONB containment: '["Bathtub", "Balcony"]' @> '["Bathtub"]'
			// GORM JSON query
			query = query.Where("room_offers.amenities @> ?", fmt.Sprintf("[\"%s\"]", amenity))
		}
	}

	if joinRooms {
		query = query.Joins("JOIN room_offers ON room_offers.hotel_id = hotels.hotel_code").Distinct()
	}

	// --- Hotel Filters ---
	if stars := c.Query("stars"); stars != "" {
		query = query.Where("star_rating >= ?", stars)
	}
	if rating := c.Query("rating"); rating != "" {
		query = query.Where("user_rating >= ?", rating)
	}
	if pType := c.Query("type"); pType != "" {
		query = query.Where("property_type ILIKE ?", pType)
	}

	// Policy Filters (JSONB key exists or value matches)
	// ?policies=alcohol:allowed,pets:allowed
	if policies := c.Query("policies"); policies != "" {
		policyMap := utils.ParseKeyVal(policies) // Implement this in utils or inline
		for key, val := range policyMap {
			// e.g. policies->>'alcohol' = 'allowed'
			query = query.Where(fmt.Sprintf("policies->>'%s' = ?", key), val)
		}
	}

	// --- Banquet Filters ---
	joinBanquets := false
	if hallType := c.Query("hall_type"); hallType != "" {
		joinBanquets = true
		query = query.Where("banquet_halls.hall_type ILIKE ?", hallType)
	}
	if hallCapacity := c.Query("hall_capacity"); hallCapacity != "" {
		joinBanquets = true
		query = query.Where("banquet_halls.capacity >= ?", hallCapacity)
	}
	if minArea := c.Query("min_hall_area"); minArea != "" {
		joinBanquets = true
		query = query.Where("banquet_halls.area >= ?", minArea)
	}
	if minHeight := c.Query("min_ceiling_height"); minHeight != "" {
		joinBanquets = true
		query = query.Where("banquet_halls.height >= ?", minHeight)
	}
	// ?banquet_features=AV,Projector
	if banquetFeatures := c.Query("banquet_features"); banquetFeatures != "" {
		joinBanquets = true
		featuresList := utils.ParseCSV(banquetFeatures)
		for _, feature := range featuresList {
			query = query.Where("banquet_halls.features @> ?", fmt.Sprintf("[\"%s\"]", feature))
		}
	}

	if joinBanquets {
		query = query.Joins("JOIN banquet_halls ON banquet_halls.hotel_id = hotels.hotel_code").Distinct()
	}

	// --- Catering Filters ---
	joinCatering := false
	if dietary := c.Query("dietary"); dietary != "" {
		joinCatering = true
		// Check if any menu has the tag
		query = query.Where("catering_menus.dietary_tags @> ?", fmt.Sprintf("[\"%s\"]", dietary))
	}

	if joinCatering {
		query = query.Joins("JOIN catering_menus ON catering_menus.hotel_id = hotels.hotel_code").Distinct()
	}

	// --- Location Tags ---
	// ?location_tags=Near Beach,City Center
	if locTags := c.Query("location_tags"); locTags != "" {
		tagsList := utils.ParseCSV(locTags)
		for _, tag := range tagsList {
			query = query.Where("location_tags @> ?", fmt.Sprintf("[\"%s\"]", tag))
		}
	}

	// --- Hotel Facilities (General) ---
	// ?facilities=Spa,Gym,Prayer Room,Kids Club
	if facilities := c.Query("facilities"); facilities != "" {
		fList := utils.ParseCSV(facilities)
		for _, f := range fList {
			query = query.Where("facilities @> ?", fmt.Sprintf("[\"%s\"]", f))
		}
	}

	// Execute Query
	result := query.Limit(50).Find(&hotels)

	if result.Error != nil {
		return utils.InternalErrorResponse(c, "Failed to fetch hotels")
	}

	// 4. Store in Redis
	if store.RDB != nil {
		if data, err := json.Marshal(hotels); err == nil {
			store.RDB.Set(ctx, cacheKey, data, 15*24*time.Hour)
		}
	}

	// 5. Return Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"count":  len(hotels),
		"data":   hotels,
	})
}

// GetHotel fetches a single hotel by ID (with full details)
// GET /api/v1/hotels/:id
func (r *Repository) GetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")

	var hotel models.Hotel
	cacheKey := "hotel:" + hotelID
	ctx := context.Background()

	// 1. Try to get from Redis
	if store.RDB != nil {
		cachedData, err := store.RDB.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedData), &hotel); err == nil {
				log.Printf("⚡ [REDIS] CACHE HIT: %s\n", cacheKey)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status": "success",
					"data":   hotel,
				})
			}
		}
	}

	// 2. Query Database
	result := store.DB.Preload("Rooms").Preload("Banquets").Preload("Menus").First(&hotel, "hotel_code = ?", hotelID)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Hotel not found",
		})
	}

	// 3. Store in Redis
	if store.RDB != nil {
		if data, err := json.Marshal(hotel); err == nil {
			store.RDB.Set(ctx, cacheKey, data, 15*24*time.Hour)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   hotel,
	})
}

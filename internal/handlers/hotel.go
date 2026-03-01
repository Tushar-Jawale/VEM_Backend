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
	"gorm.io/gorm"
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
	skipCache := c.Query("skip_cache") == "true"
	if store.RDB != nil && !skipCache {
		cachedData, err := store.RDB.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedData), &hotels); err == nil {

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
		Preload("Banquets").Preload("Menus")

	// --- Room Filters ---
	// Consolidate all room-level constraints into a shared subquery builder
	// to ensure strict AND logic (Hotel must have rooms that satisfy EVERY filter)

	minPrice := utils.ParseFloat64(c.Query("min_price"), 0)
	maxPrice := utils.ParseFloat64(c.Query("max_price"), 0)

	freeCancellation := c.Query("free_cancellation") == "true"

	guestsPerRoom := c.Query("guests_per_room")
	if guestsPerRoom == "" {
		guestsPerRoom = c.Query("occupancy")
	}
	gValue := utils.ParseInt(guestsPerRoom, 0)

	roomAmenities := c.Query("room_amenities")
	amenitiesList := utils.ParseCSV(roomAmenities)

	// Build the base EXISTS subquery fragment
	// exactCap=true → use max_capacity = ? (strict room-type match)
	// exactCap=false → use max_capacity >= ? (general occupancy / "fits at least N people")
	buildRoomSubQuery := func(alias string, minCap int, minCount int, exactCap bool) (string, []interface{}) {
		sql := fmt.Sprintf("SELECT 1 FROM room_offers %s WHERE %s.hotel_id = hotels.hotel_code", alias, alias)
		var args []interface{}

		if minCap > 0 {
			if exactCap {
				sql += fmt.Sprintf(" AND %s.max_capacity = ?", alias)
			} else {
				sql += fmt.Sprintf(" AND %s.max_capacity >= ?", alias)
			}
			args = append(args, minCap)
		}
		if minCount > 0 {
			sql += fmt.Sprintf(" AND %s.count >= ?", alias)
			args = append(args, minCount)
		}
		if minPrice > 0 {
			sql += fmt.Sprintf(" AND %s.total_fare >= ?", alias)
			args = append(args, minPrice)
		}
		if maxPrice > 0 {
			sql += fmt.Sprintf(" AND %s.total_fare <= ?", alias)
			args = append(args, maxPrice)
		}
		if freeCancellation {
			sql += fmt.Sprintf(" AND %s.is_refundable = ?", alias)
			args = append(args, true)
		}
		for _, amenity := range amenitiesList {
			sql += fmt.Sprintf(" AND %s.amenities @> ?", alias)
			args = append(args, fmt.Sprintf("[\"%s\"]", amenity))
		}
		return sql, args
	}

	// 1. Room Configuration (Complex Inventory Filter)
	type Req struct {
		Occupancy int `json:"occupancy"`
		Count     int `json:"count"`
	}
	var reqs []Req
	if roomConfig := c.Query("room_config"); roomConfig != "" {
		if err := json.Unmarshal([]byte(roomConfig), &reqs); err != nil {
			log.Printf("⚠️  Failed to unmarshal room_config: %v", err)
		}
	}

	// Specific room type requirements — use EXACT capacity matching
	if s := c.QueryInt("rooms_single", 0); s > 0 {
		reqs = append(reqs, Req{Occupancy: 1, Count: s})
	}
	if d := c.QueryInt("rooms_double", 0); d > 0 {
		reqs = append(reqs, Req{Occupancy: 2, Count: d})
	}
	if t := c.QueryInt("rooms_triple", 0); t > 0 {
		reqs = append(reqs, Req{Occupancy: 3, Count: t})
	}
	if q := c.QueryInt("rooms_quad", 0); q > 0 {
		reqs = append(reqs, Req{Occupancy: 4, Count: q})
	}

	if len(reqs) > 0 {
		for i, r := range reqs {
			// Use exact capacity matching for named room types
			sql, args := buildRoomSubQuery(fmt.Sprintf("ro%d", i), r.Occupancy, r.Count, true)

			query = query.Where(fmt.Sprintf("EXISTS (%s)", sql), args...)
		}
	} else if gValue > 0 || minPrice > 0 || maxPrice > 0 || freeCancellation || len(amenitiesList) > 0 {
		// General occupancy filter: use >= (any room that fits at least N people)
		sql, args := buildRoomSubQuery("ro_base", gValue, 0, false)
		query = query.Where(fmt.Sprintf("EXISTS (%s)", sql), args...)
	}

	// Filter preloaded rooms so the user sees exactly what they searched for
	query = query.Preload("Rooms", func(db *gorm.DB) *gorm.DB {
		preloadQuery := db
		if len(reqs) > 0 {
			// Only preload rooms whose capacity exactly matches one of the requested types
			var caps []int
			for _, r := range reqs {
				caps = append(caps, r.Occupancy)
			}

			preloadQuery = preloadQuery.Where("max_capacity IN ?", caps)
		} else if gValue > 0 {
			// General occupancy: any room that can fit at least N people
			preloadQuery = preloadQuery.Where("max_capacity >= ?", gValue)
		}
		if minPrice > 0 {
			preloadQuery = preloadQuery.Where("total_fare >= ?", minPrice)
		}
		if maxPrice > 0 {
			preloadQuery = preloadQuery.Where("total_fare <= ?", maxPrice)
		}
		if freeCancellation {
			preloadQuery = preloadQuery.Where("is_refundable = ?", true)
		}
		for _, amenity := range amenitiesList {
			preloadQuery = preloadQuery.Where("amenities @> ?", fmt.Sprintf("[\"%s\"]", amenity))
		}
		return preloadQuery
	})

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

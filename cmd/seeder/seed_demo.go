package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/akashtripathi12/TBO_Backend/internal/models"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
	"github.com/google/uuid"
)

func SeedDemoData() {
	log.Println("🚀 Starting Specialized Demo Seeding...")

	// 1. Target Countries (5)
	targetCountries := []string{"IN", "US", "SG", "AE", "TH"}
	SeedCountries(targetCountries)

	// 2. Cities (50 total, 10 per country)
	// We'll use the existing SeedCities but we need to ensure we only get a certain amount
	// For simplicity in this specialized script, we'll fetch them from the database after seeding all
	SeedCities(targetCountries)

	var allCities []models.City
	store.DB.Find(&allCities)

	cityMap := make(map[string][]models.City)
	for _, city := range allCities {
		cityMap[city.CountryCode] = append(cityMap[city.CountryCode], city)
	}

	// 3. Selection of 10 cities per country
	var demoCities []models.City
	for _, countryCode := range targetCountries {
		cities := cityMap[countryCode]
		if len(cities) > 10 {
			cities = cities[:10]
		}
		demoCities = append(demoCities, cities...)
	}

	log.Printf("🏙️  Selected %d cities for demo seeding.", len(demoCities))

	// 4. Hotels (100 total)
	// 50 for India, 50 for others (12-13 each for US, SG, AE, TH)
	hotalsPerCountry := make(map[string]int)
	hotalsPerCountry["IN"] = 50
	hotalsPerCountry["US"] = 13
	hotalsPerCountry["SG"] = 12
	hotalsPerCountry["AE"] = 13
	hotalsPerCountry["TH"] = 12

	totalHotelsSeeded := 0

	for countryCode, count := range hotalsPerCountry {
		cities := cityMap[countryCode]
		if len(cities) == 0 {
			continue
		}

		hotelsPerCity := count / len(cities)
		if hotelsPerCity == 0 {
			hotelsPerCity = 1
		}

		seededForCountry := 0
		for _, city := range cities {
			if seededForCountry >= count {
				break
			}

			// Fetch hotels from API for this city
			hotels, err := fetchHotelsForCity(city.ID)
			if err != nil || len(hotels) == 0 {
				continue
			}

			limit := hotelsPerCity
			if seededForCountry+limit > count {
				limit = count - seededForCountry
			}
			if limit > len(hotels) {
				limit = len(hotels)
			}

			for i := 0; i < limit; i++ {
				h := hotels[i]
				dbHotel := models.Hotel{
					ID:         h.HotelCode,
					CityID:     city.ID,
					Name:       h.HotelName,
					StarRating: convertStarRating(h.HotelRating),
					Address:    h.Address,
				}
				store.DB.FirstOrCreate(&dbHotel)

				// 5. Seed Rooms for this hotel (5 types)
				seedDemoRooms(dbHotel.ID)

				// 6. Seed Banquets and Catering for this hotel
				seedDemoBanquetsAndCatering(dbHotel.ID)

				seededForCountry++
				totalHotelsSeeded++
			}
		}
		log.Printf("🏨 Seeded %d hotels for %s", seededForCountry, countryCode)
	}

	log.Printf("🎉 Demo Seeding Completed! Seeded %d hotels with rooms, banquets, and catering.", totalHotelsSeeded)
}

func seedDemoRooms(hotelID string) {
	// Capacities 1, 2, 3, 4 and one more
	capacities := []int{1, 2, 3, 4, 2} // repeating 2 or can be random 5+
	roomNames := []string{"Single Deluxe", "Double Premium", "Triple Suite", "Family Quad", "Executive King"}

	for i := 0; i < 5; i++ {
		room := models.RoomOffer{
			ID:          uuid.New().String(),
			HotelID:     hotelID,
			Name:        roomNames[i],
			MaxCapacity: capacities[i],
			Count:       rand.Intn(11) + 90, // 90 to 100
			BookingCode: fmt.Sprintf("DEMO-%s-%d", hotelID, i),
			TotalFare:   float64((rand.Intn(100) + 50) * 10),
			Currency:    "USD",
		}
		store.DB.Create(&room)
	}
}

func seedDemoBanquetsAndCatering(hotelID string) {
	// 4-5 Banquets
	numBanquets := rand.Intn(2) + 4
	for i := 0; i < numBanquets; i++ {
		name := banquetNames[rand.Intn(len(banquetNames))]
		if numBanquets > 1 {
			name = fmt.Sprintf("%s %d", name, i+1)
		}
		banquet := models.BanquetHall{
			HotelID:     hotelID,
			Name:        name,
			Capacity:    (rand.Intn(10) + 5) * 50,
			PricePerDay: float64((rand.Intn(20) + 10) * 100),
		}
		store.DB.Create(&banquet)
	}

	// 4-5 Menus
	numMenus := rand.Intn(2) + 4
	for i := 0; i < numMenus; i++ {
		name := cateringMenuNames[rand.Intn(len(cateringMenuNames))]
		menuType := "veg"
		if rand.Intn(2) == 0 {
			menuType = "non-veg"
		}
		menu := models.CateringMenu{
			HotelID:       hotelID,
			Name:          name,
			Type:          menuType,
			PricePerPlate: float64((rand.Intn(5) + 5) * 10),
		}
		store.DB.Create(&menu)
	}
}

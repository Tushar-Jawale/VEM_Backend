package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/akashtripathi12/TBO_Backend/internal/models"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
)

var banquetNames = []string{
	"Grand Ballroom",
	"Crystal Room",
	"Royal Hall",
	"Sky Lounge",
	"Garden Pavilion",
	"Emerald Suite",
	"Sapphire Lounge",
	"Diamond Gallery",
	"Majestic Venue",
	"Zen Terrace",
}

var cateringMenuNames = []string{
	"Executive Lunch",
	"Royal Wedding Feast",
	"Classic High Tea",
	"International Buffet",
	"Traditional Asian Platter",
	"Gourmet Italian Night",
	"Seafood Special",
	"Vegan Delight",
	"Continental Breakfast",
	"Midnight Snacks",
}

func SeedBanquets() {
	log.Println("🍽️  Seeding Banquet Halls and Catering Menus...")

	var hotels []models.Hotel
	if err := store.DB.Find(&hotels).Error; err != nil {
		log.Fatalf("❌ Failed to fetch hotels: %v", err)
	}

	totalBanquets := 0
	totalMenus := 0

	for _, hotel := range hotels {
		// Seed 4-5 Banquets
		numBanquets := rand.Intn(2) + 4 // 4 or 5
		for i := 0; i < numBanquets; i++ {
			name := banquetNames[rand.Intn(len(banquetNames))]
			if numBanquets > 1 {
				name = fmt.Sprintf("%s %d", name, i+1)
			}
			banquet := models.BanquetHall{
				HotelID:     hotel.ID,
				Name:        name,
				Capacity:    (rand.Intn(10) + 5) * 50,            // 250 to 700
				PricePerDay: float64((rand.Intn(20) + 10) * 100), // 1000 to 3000
			}
			store.DB.Create(&banquet)
			totalBanquets++
		}

		// Seed 4-5 Menus
		numMenus := rand.Intn(2) + 4 // 4 or 5
		for i := 0; i < numMenus; i++ {
			name := cateringMenuNames[rand.Intn(len(cateringMenuNames))]
			menuType := "veg"
			if rand.Intn(2) == 0 {
				menuType = "non-veg"
			}
			menu := models.CateringMenu{
				HotelID:       hotel.ID,
				Name:          name,
				Type:          menuType,
				PricePerPlate: float64((rand.Intn(5) + 5) * 10), // 50 to 100
			}
			store.DB.Create(&menu)
			totalMenus++
		}
	}

	log.Printf("✅ Seeded %d Banquet Halls and %d Catering Menus across %d hotels.", totalBanquets, totalMenus, len(hotels))
}

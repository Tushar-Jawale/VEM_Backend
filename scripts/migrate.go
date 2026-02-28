package main

import (
	"log"

	"github.com/akashtripathi12/TBO_Backend/internal/config"
	"github.com/akashtripathi12/TBO_Backend/internal/models"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")
	config.Load()
	store.InitDB()

	log.Println("Running AutoMigrate...")
	err := store.DB.AutoMigrate(&models.Event{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration successful!")
}

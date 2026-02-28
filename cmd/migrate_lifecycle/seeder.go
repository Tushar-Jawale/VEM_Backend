package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env") // fallback
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	// Connect to DB directly using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("🚀 Starting Lifecycle Migration...")

	// 1. Normalize existing data
	// UPDATE events SET status = 'active' WHERE status IN ('draft', 'allocating');
	if err := db.Exec("UPDATE events SET status = 'active' WHERE status IN ('draft', 'allocating')").Error; err != nil {
		log.Fatalf("❌ Failed to update draft/allocating events: %v", err)
	}
	log.Println("✅ Updated legacy statuses (draft, allocating) -> active")

	// UPDATE events SET status = 'finalized' WHERE status = 'rooms_finalized';
	if err := db.Exec("UPDATE events SET status = 'finalized' WHERE status = 'rooms_finalized'").Error; err != nil {
		log.Fatalf("❌ Failed to update rooms_finalized events: %v", err)
	}
	log.Println("✅ Updated legacy statuses (rooms_finalized) -> finalized")

	// 2. Drop old constraint if exists
	// ALTER TABLE events DROP CONSTRAINT IF EXISTS valid_status;
	if err := db.Exec("ALTER TABLE events DROP CONSTRAINT IF EXISTS valid_status").Error; err != nil {
		log.Fatalf("❌ Failed to drop old constraint: %v", err)
	}
	log.Println("✅ Dropped old 'valid_status' constraint")

	// 3. Add new constraint
	// ALTER TABLE events ADD CONSTRAINT valid_status CHECK (status IN ('active', 'finalized'));
	if err := db.Exec("ALTER TABLE events ADD CONSTRAINT valid_status CHECK (status IN ('active', 'finalized'))").Error; err != nil {
		log.Fatalf("❌ Failed to add new constraint: %v", err)
	}
	log.Println("✅ Added new 'valid_status' constraint ('active', 'finalized')")

	log.Println("🎉 Migration Completed Successfully!")
}

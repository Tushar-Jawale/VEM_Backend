package store

import (
	"context"
	"log"

	"github.com/akashtripathi12/TBO_Backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ctx = context.Background()

func InitRedis(cfg *config.Config) {
	if cfg.RedisAddr == "" {
		log.Println("⚠️ Redis address not provided, skipping Redis initialization")
		return
	}

	// --- FIX: Parse the full Redis URL instead of manually setting Options ---
	opt, err := redis.ParseURL(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("❌ Failed to parse Redis URL: %v", err)
	}

	// Initialize using the parsed options
	RDB = redis.NewClient(opt)

	// Test the connection
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		log.Println("❌ Failed to connect to Redis:", err)
	} else {
		log.Println("✅ Connected to Redis at", cfg.RedisAddr)

		// Set LFU strategy
		err = RDB.ConfigSet(ctx, "maxmemory-policy", "allkeys-lfu").Err()
		if err != nil {
			log.Println("⚠️ Failed to set Redis eviction policy to LFU:", err)
		} else {
			log.Println("✅ Redis eviction policy set to allkeys-lfu")
		}
	}
}

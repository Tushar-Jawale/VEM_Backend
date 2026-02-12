package main

import (
	"log"
	"net/http"

	"github.com/akashtripathi12/TBO_Backend/internal/config"
	"github.com/akashtripathi12/TBO_Backend/internal/handlers"
	"github.com/akashtripathi12/TBO_Backend/internal/routes"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
)

func main() {
	// 1. Load Configuration
	appConfig := config.Load()
	log.Printf("Loaded configuration for env: %s", appConfig.Env)

	// 2. Initialize Store
	db := store.NewMockStore()

	// 3. Initialize Repository/Handlers
	repo := handlers.NewRepository(appConfig, db)

	// 4. Setup Routes
	srv := &http.Server{
		Addr:    appConfig.Port,
		Handler: routes.Routes(appConfig, repo),
	}

	// 5. Start Server
	log.Printf("Starting server on %s...", appConfig.Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

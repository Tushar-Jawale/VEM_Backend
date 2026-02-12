package handlers

import (
	"github.com/akashtripathi12/TBO_Backend/internal/config"
	"github.com/akashtripathi12/TBO_Backend/internal/store"
)

// Repository holds the application configuration and database store
type Repository struct {
	App *config.Config
	DB  store.Store
}

// NewRepository creates a new instance of the repository
func NewRepository(app *config.Config, db store.Store) *Repository {
	return &Repository{
		App: app,
		DB:  db,
	}
}

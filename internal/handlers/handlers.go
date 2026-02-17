package handlers

import (
	"github.com/akashtripathi12/TBO_Backend/internal/config"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type Repository struct {
	App         *config.Config
	DB          *gorm.DB
	QueueClient *asynq.Client
}

// NewRepository creates a new instance of the repository
func NewRepository(app *config.Config, db *gorm.DB, queueClient *asynq.Client) *Repository {
	return &Repository{
		App:         app,
		DB:          db,
		QueueClient: queueClient,
	}
}

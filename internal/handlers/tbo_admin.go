package handlers

import (
	"github.com/akashtripathi12/TBO_Backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// ListNegotiations returns all active and locked negotiation sessions.
// Intended for the TBO Agent dashboard.
func (r *Repository) ListNegotiations(c *fiber.Ctx) error {
	// For a real app, we should enforce that the caller has Role == "tbo_agent"
	// userRole := c.Locals("userRole").(string)
	// if userRole != "tbo_agent" {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	// }

	var sessions []models.NegotiationSession
	// Preload the associated Event to show details on the frontend dashboard
	if err := r.DB.Preload("Event").Order("updated_at desc").Find(&sessions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch negotiation sessions"})
	}

	return c.JSON(sessions)
}

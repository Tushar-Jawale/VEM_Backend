package routes

import (
	"net/http"

	"github.com/akashtripathi12/TBO_Backend/internal/config"
	"github.com/akashtripathi12/TBO_Backend/internal/handlers"
	"github.com/akashtripathi12/TBO_Backend/internal/middleware"
)

func Routes(app *config.Config, repo *handlers.Repository) http.Handler {
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("POST /api/v1/auth/login/agent", repo.LoginAgent)
	mux.HandleFunc("POST /api/v1/auth/login/guest", repo.LoginGuest)
	mux.HandleFunc("POST /api/v1/auth/logout", repo.Logout)
	mux.HandleFunc("GET /api/v1/auth/me", repo.GetCurrentUser)

	// Events
	mux.HandleFunc("GET /api/v1/events", repo.GetEvents)
	mux.HandleFunc("POST /api/v1/events", repo.CreateEvent)
	mux.HandleFunc("GET /api/v1/events/{id}", repo.GetEvent)
	mux.HandleFunc("PUT /api/v1/events/{id}", repo.UpdateEvent)
	mux.HandleFunc("DELETE /api/v1/events/{id}", repo.DeleteEvent)
	mux.HandleFunc("GET /api/v1/dashboard/metrics", repo.GetMetrics)
	mux.HandleFunc("GET /api/v1/events/{id}/venues", repo.GetEventVenues)
	mux.HandleFunc("GET /api/v1/events/{id}/allocations", repo.GetEventAllocations)

	// Guests
	mux.HandleFunc("GET /api/v1/events/{id}/guests", repo.GetGuests)
	mux.HandleFunc("GET /api/v1/guests/{id}", repo.GetGuest)
	mux.HandleFunc("PUT /api/v1/guests/{id}", repo.UpdateGuest)
	mux.HandleFunc("DELETE /api/v1/guests/{id}", repo.DeleteGuest)
	mux.HandleFunc("POST /api/v1/guests", repo.CreateGuest)
	mux.HandleFunc("POST /api/v1/guests/{id}/subguests", repo.AddSubGuest)

	// Allocations
	mux.HandleFunc("POST /api/v1/allocations", repo.CreateAllocation)
	mux.HandleFunc("PUT /api/v1/allocations/{id}", repo.UpdateAllocation)

	return middleware.EnableCORS(mux)
}

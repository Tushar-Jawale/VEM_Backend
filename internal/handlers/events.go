package handlers

import "net/http"

func (m *Repository) GetEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Get events from store
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Events Endpoint"))
}

func (m *Repository) GetEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Get event by ID
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Event Endpoint"))
}

func (m *Repository) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Create new event
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create Event Endpoint"))
}

func (m *Repository) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Get metrics
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Metrics Endpoint"))
}

func (m *Repository) GetEventVenues(w http.ResponseWriter, r *http.Request) {
	// TODO: Get venues
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Event Venues Endpoint"))
}

func (m *Repository) GetEventAllocations(w http.ResponseWriter, r *http.Request) {
	// TODO: Get allocations
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Event Allocations Endpoint"))
}
func (m *Repository) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Update event
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update Event Endpoint"))
}

func (m *Repository) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete event
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete Event Endpoint"))
}

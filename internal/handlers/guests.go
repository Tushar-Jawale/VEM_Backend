package handlers

import "net/http"

func (m *Repository) GetGuests(w http.ResponseWriter, r *http.Request) {
	// TODO: Get guests by event ID
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Guests Endpoint"))
}

func (m *Repository) GetGuest(w http.ResponseWriter, r *http.Request) {
	// TODO: Get guest by ID
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Guest Endpoint"))
}

func (m *Repository) CreateGuest(w http.ResponseWriter, r *http.Request) {
	// TODO: Create guest
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create Guest Endpoint"))
}

func (m *Repository) AddSubGuest(w http.ResponseWriter, r *http.Request) {
	// TODO: Add sub guest
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Add Sub Guest Endpoint"))
}
func (m *Repository) UpdateGuest(w http.ResponseWriter, r *http.Request) {
	// TODO: Update guest
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update Guest Endpoint"))
}

func (m *Repository) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	// TODO: Delete guest
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete Guest Endpoint"))
}

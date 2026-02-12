package handlers

import "net/http"

func (m *Repository) CreateAllocation(w http.ResponseWriter, r *http.Request) {
	// TODO: Create allocation
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create Allocation Endpoint"))
}

func (m *Repository) UpdateAllocation(w http.ResponseWriter, r *http.Request) {
	// TODO: Update allocation
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update Allocation Endpoint"))
}

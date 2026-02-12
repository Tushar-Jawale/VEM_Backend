package handlers

import "net/http"

func (m *Repository) LoginAgent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement agent login logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Agent Login Endpoint"))
}

func (m *Repository) LoginGuest(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement guest login logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Guest Login Endpoint"))
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logout logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout Endpoint"))
}

func (m *Repository) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get current user logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get Current User Endpoint"))
}

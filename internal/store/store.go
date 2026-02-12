package store

import "github.com/akashtripathi12/TBO_Backend/internal/models"

type Store interface {
	// Auth
	GetUserByEmail(email string) (*models.User, error)
	GetAgentCredentials() (*models.AuthCredentials, error)

	// Events
	GetEvents() ([]models.Event, error)
	GetEventByID(id string) (*models.Event, error)
	CreateEvent(event models.Event) error
	GetMetrics() ([]models.MetricData, error)

	// Guests
	GetGuestsByEventID(eventID string) ([]models.HeadGuest, error)
	GetGuestByID(id string) (*models.HeadGuest, error)
	AddHeadGuest(guest models.HeadGuest) error

	// SubGuests
	GetSubGuestsByHeadGuestID(headGuestID string) ([]models.SubGuest, error)

	// Allocations
	GetAllocationsByEventID(eventID string) ([]models.RoomAllocation, error)

	// Venues
	GetVenuesByEventID(eventID string) ([]models.CuratedVenue, error)
}

type MockStore struct{}

func NewMockStore() *MockStore {
	return &MockStore{}
}

// Implement Store interface methods (returning nil/empty for now to satisfy interface)
func (m *MockStore) GetUserByEmail(email string) (*models.User, error)             { return nil, nil }
func (m *MockStore) GetAgentCredentials() (*models.AuthCredentials, error)         { return nil, nil }
func (m *MockStore) GetEvents() ([]models.Event, error)                            { return nil, nil }
func (m *MockStore) GetEventByID(id string) (*models.Event, error)                 { return nil, nil }
func (m *MockStore) CreateEvent(event models.Event) error                          { return nil }
func (m *MockStore) GetMetrics() ([]models.MetricData, error)                      { return nil, nil }
func (m *MockStore) GetGuestsByEventID(eventID string) ([]models.HeadGuest, error) { return nil, nil }
func (m *MockStore) GetGuestByID(id string) (*models.HeadGuest, error)             { return nil, nil }
func (m *MockStore) AddHeadGuest(guest models.HeadGuest) error                     { return nil }
func (m *MockStore) GetSubGuestsByHeadGuestID(headGuestID string) ([]models.SubGuest, error) {
	return nil, nil
}
func (m *MockStore) GetAllocationsByEventID(eventID string) ([]models.RoomAllocation, error) {
	return nil, nil
}
func (m *MockStore) GetVenuesByEventID(eventID string) ([]models.CuratedVenue, error) {
	return nil, nil
}

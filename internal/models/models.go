package models

// Event represents a wedding or corporate event
type Event struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Location          string `json:"location"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	Organizer         string `json:"organizer"`
	GuestCount        int    `json:"guestCount"`
	HotelCount        int    `json:"hotelCount"`
	InventoryConsumed int    `json:"inventoryConsumed"`
	Status            string `json:"status"` // 'active' | 'upcoming' | 'completed'
}

// HeadGuest represents the primary contact for a group of guests
type HeadGuest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	EventID      string `json:"eventId"`
	SubGroupName string `json:"subGroupName"` // e.g., "Bride's Family"
}

// SubGuest represents a guest under a HeadGuest
type SubGuest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Age         int    `json:"age,omitempty"`
	GuestCount  int    `json:"guestCount"` // Defaults to 1
	HeadGuestID string `json:"headGuestId"`
	RoomGroupID string `json:"roomGroupId,omitempty"` // null if unassigned
}

// RoomAllocation represents a block of rooms allocated to a HeadGuest
type RoomAllocation struct {
	ID          string `json:"id"`
	EventID     string `json:"eventId"`
	HeadGuestID string `json:"headGuestId"`
	RoomType    string `json:"roomType"` // e.g., "Deluxe Room", "Suite"
	MaxCapacity int    `json:"maxCapacity"`
	HotelName   string `json:"hotelName"`
}

// RoomGroup represents a specific assignment of guests to a room
type RoomGroup struct {
	ID           string   `json:"id"`
	AllocationID string   `json:"allocationId"`
	GuestIDs     []string `json:"guestIds"` // SubGuest IDs assigned to this room
	CustomLabel  string   `json:"customLabel,omitempty"`
}

// CuratedVenue represents a hotel or venue option
type CuratedVenue struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Amenities   []string `json:"amenities"`
	EventID     string   `json:"eventId"`
}

// MetricData represents dashboard analytics
type MetricData struct {
	Label  string      `json:"label"`
	Value  interface{} `json:"value"` // string or number
	Change int         `json:"change,omitempty"`
	Trend  string      `json:"trend,omitempty"` // 'up' | 'down' | 'neutral'
}

// User represents an authenticated user (Agent or Guest)
type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"` // 'agent' | 'guest'
	EventID string `json:"eventId,omitempty"`
}

// AuthCredentials for login
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

package trips

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID           uuid.UUID  `json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	Status       TripStatus `json:"status"`
	Name         string     `json:"name"`
	Destination  string     `json:"destination"`
	Origin       string     `json:"origin"`
	DateFrom     Date       `json:"date_from"`
	DateTo       Date       `json:"date_to"`
	Budget       int64      `json:"budget"`
	Requirements string     `json:"requirements"`
}

type TripData struct {
	Name         string
	Destination  string
	Origin       string
	DateFrom     Date `json:"date_from"`
	DateTo       Date `json:"date_to"`
	Budget       int64
	Requirements string
}

type TripResult struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	TripID      uuid.UUID
	Summary     string
	Attractions string
	Weather     string
	Prices      string
	Luggage     string
	Documents   string
	Commuting   string
}

package trips

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Name         string
	Destination  string
	Origin       string
	DateFrom     time.Time
	DateTo       time.Time
	Budget       int64
	Requirements string
}

type TripData struct {
	Name         string
	Destination  string
	Origin       string
	DateFrom     time.Time
	DateTo       time.Time
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

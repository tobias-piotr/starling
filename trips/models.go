package trips

import (
	"time"

	"github.com/google/uuid"
)

type TripStatus int64

const (
	Draft     TripStatus = iota
	Requested            = iota
	Failed               = iota
	Completed            = iota
)

func (s TripStatus) String() string {
	return [...]string{"draft", "requested", "failed", "completed"}[s]
}

type Trip struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	Status       TripStatus
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

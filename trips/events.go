package trips

import "github.com/google/uuid"

// TripCreated is an event that is published when a trip is created
type TripCreated struct {
	tripID uuid.UUID
}

func (e TripCreated) String() string {
	return "trip_created"
}

func (e TripCreated) Payload() map[string]any {
	return map[string]any{"trip_id": e.tripID}
}

// TripRequested is an event that is published when a trip is requested
type TripRequested struct {
	tripID uuid.UUID
}

func (e TripRequested) String() string {
	return "trip_requested"
}

func (e TripRequested) Payload() map[string]any {
	return map[string]any{"trip_id": e.tripID}
}

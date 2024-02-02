package trips

import (
	"starling/internal/events"
)

// TripService is an orchestrator for the trip domain
type TripService struct {
	tripRepository TripRepository
	eventBus       events.EventBus
}

func (s *TripService) CreateTrip(data *TripData) (*Trip, error) {
	// Create trip
	trip, err := s.tripRepository.Create(data)
	if err != nil {
		return nil, err
	}

	// Publish event
	go s.eventBus.Publish(TripCreated.String(), map[string]any{"trip_id": trip.ID})

	return trip, nil
}

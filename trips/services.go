package trips

import (
	"starling/internal/events"
)

// TripService is an orchestrator for the trip domain
type TripService struct {
	tripRepository TripRepository
	eventBus       events.EventBus
}

func NewTripService(tripRepository TripRepository, eventBus events.EventBus) *TripService {
	return &TripService{tripRepository, eventBus}
}

func (s *TripService) CreateTrip(data *TripData) (*Trip, error) {
	// Create trip
	trip, err := s.tripRepository.Create(data)
	if err != nil {
		return nil, err
	}

	// Publish event
	go s.eventBus.Publish(TripCreated{tripID: trip.ID})

	return trip, nil
}

func (s *TripService) GetTrips() ([]*Trip, error) {
	// TODO: Add pagination
	return s.tripRepository.GetAll()
}

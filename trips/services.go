package trips

import (
	"log/slog"

	"starling/internal/domain"
	"starling/internal/events"
)

// TripService is an orchestrator for the trip domain.
type TripService struct {
	tripRepository TripRepository
	eventBus       events.EventBus
}

func NewTripService(tripRepository TripRepository, eventBus events.EventBus) *TripService {
	return &TripService{tripRepository, eventBus}
}

func (s *TripService) CreateTrip(data *TripData) (*Trip, error) {
	slog.Info("Creating a new trip", "data", data)

	// Create trip
	trip, err := s.tripRepository.Create(data)
	if err != nil {
		return nil, err
	}

	// Publish event
	go s.eventBus.Publish(TripCreated{tripID: trip.ID})

	return trip, nil
}

func (s *TripService) GetTrips(page int, perPage int) ([]*TripOverview, error) {
	return s.tripRepository.GetAll(page, perPage)
}

func (s *TripService) GetTrip(id string) (*Trip, error) {
	return s.tripRepository.Get(id)
}

func (s *TripService) RequestTrip(id string) error {
	slog.Info("Requesting trip", "id", id)

	trip, err := s.tripRepository.Get(id)
	if err != nil {
		return err
	}
	if trip == nil {
		return &domain.NotFoundErr{Msg: "Trip not found"}
	}

	if err := trip.ValidateRequest(); err != nil {
		return &domain.ValidationErr{Err: err}
	}

	if err := s.tripRepository.Update(
		id,
		map[string]any{"status": RequestedStatus.String()},
	); err != nil {
		return err
	}

	go s.eventBus.Publish(TripRequested{tripID: trip.ID})

	return nil
}

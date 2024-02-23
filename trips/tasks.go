package trips

import (
	"fmt"
	"log/slog"
	"sync"

	"starling/internal/ai"
	"starling/internal/domain"
)

const (
	basePrompt = `
You are a travel asssistant. You help people plan trips.
People give you a set of criteria, and you give them a summary of their trip,
including aspects like weather, prices or attractions.
Don't reply with everything at once. Questions about specific sections will
be asked separately.
`

	summaryPrompt = `
Generate a summary for the following criteria:
Origin: %s
Destination: %s
Date: %s to %s
Budget: %d
Requirements: %s
Summary should be a general, few sentences long description of the destination.
Talk only about the place itself, not the trip, weather, prices or attractions.
`
)

// RequestTrip generates an entire trip result for given trip request.
func RequestTrip(tripsRepository TripRepository, aiClient ai.AIClient, tripID string) error {
	// Get the trip
	trip, err := tripsRepository.Get(tripID)
	if err != nil {
		return fmt.Errorf("get trip: %w", err)
	}

	// Make sure the trip has a result
	if trip.Result == nil {
		res, err := tripsRepository.AddResult(tripID)
		if err != nil {
			return fmt.Errorf("add result: %w", err)
		}
		trip.Result = res
	}

	var wg sync.WaitGroup
	errs := make(chan error)

	if trip.Result.Summary == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Info("Generating summary for trip", "trip_id", trip.ID)

			summary, err := generateSummary(aiClient, trip)
			if err != nil {
				errs <- fmt.Errorf("generate summary: %w", err)
				return
			}

			if err = tripsRepository.UpdateResult(
				trip.Result.ID.String(),
				map[string]any{"summary": summary},
			); err != nil {
				errs <- fmt.Errorf("update result summary: %w", err)
				return
			}
		}()
	}

	if trip.Result.Attractions == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- fmt.Errorf("generate attractions: not implemented")
		}()
	}

	// Wait for all goroutines to finish and close the errs channel
	go func() {
		wg.Wait()
		close(errs)
	}()

	// Collect all errors
	cerr := &domain.CompositeErr{}
	for err := range errs {
		cerr.AddError(err)
	}
	if !cerr.IsEmpty() {
		err = tripsRepository.Update(
			tripID,
			map[string]any{"status": FailedStatus.String()},
		)
		if err != nil {
			return fmt.Errorf("update trip status: %w", err)
		}
		return cerr
	}

	err = tripsRepository.Update(
		tripID,
		map[string]any{"status": CompletedStatus.String()},
	)
	if err != nil {
		return fmt.Errorf("update trip status: %w", err)
	}

	return nil
}

func generateSummary(aiClient ai.AIClient, trip *Trip) (string, error) {
	prompt := fmt.Sprintf(
		basePrompt+summaryPrompt,
		trip.Origin,
		trip.Destination,
		trip.DateFrom.Format("2006-01-02"),
		trip.DateTo.Format("2006-01-02"),
		trip.Budget,
		trip.Requirements,
	)
	return aiClient.Send(prompt)
}

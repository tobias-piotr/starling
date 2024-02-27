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
Destination: %s
Date: %s to %s
Summary should be a general, few sentences long description of the destination.
Talk only about the place itself, not the trip, weather, prices or attractions.
`

	attractionsPrompt = `
Generate a list of attractions for the following criteria:
Destination: %s
Date: %s to %s
Budget: %d
Requirements: %s
Attractions should be a description of the most interesting places to visit, things to see
and activities to do in the destination.
`

	weatherPrompt = `
Generate a weather description for the following criteria:
Destination: %s
Date: %s to %s
Requirements: %s
Weather description should include the average temperature, humidity and
precipitation for the destination during the given dates. You should also consider
additional requirements that user gave, to make the description more accurate.
`

	pricesPrompt = `
Generate a price description for the following criteria:
Destination: %s
Date: %s to %s
Budget: %d
Requirements: %s
Price description should include the average price of a hotel, meal, commuting and
all the other things that user would be interested based on given criteria.
`

	luggagePrompt = `
Generate a luggage description for the following criteria:
Destination: %s
Date: %s to %s
Requirements: %s
Luggage description should include the most important things that user should
pack for the trip, based on the given criteria.
`

	documentsPrompt = `
Generate a documents description for the following criteria:
Origin: %s
Destination: %s
Date: %s to %s
Requirements: %s
Documents description should include the most important documents that user should
have for the trip, based on the given criteria.
	`

	commutingPrompt = `
Generate a commuting description for the following criteria:
Destination: %s
Date: %s to %s
Budget: %d
Requirements: %s
Commuting description should include the most convenient transportation methods
on the given destination, based on the given criteria.
`
)

var requestResultConf = []FieldConfig{
	{
		Name:   "summary",
		Prompt: summaryPrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02")}
		},
		Condition: func(t *Trip) bool { return t.Result.Summary == "" },
	},
	{
		Name:   "attractions",
		Prompt: attractionsPrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02"), t.Budget, t.Requirements}
		},
		Condition: func(t *Trip) bool { return t.Result.Attractions == "" },
	},
	{
		Name:   "weather",
		Prompt: weatherPrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02"), t.Requirements}
		},
		Condition: func(t *Trip) bool { return t.Result.Weather == "" },
	},
	{
		Name:   "prices",
		Prompt: pricesPrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02"), t.Budget, t.Requirements}
		},
		Condition: func(t *Trip) bool { return t.Result.Prices == "" },
	},
	{
		Name:   "luggage",
		Prompt: luggagePrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02"), t.Requirements}
		},
		Condition: func(t *Trip) bool { return t.Result.Luggage == "" },
	},
	{
		Name:   "documents",
		Prompt: documentsPrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Origin, t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02"), t.Requirements}
		},
		Condition: func(t *Trip) bool { return t.Result.Documents == "" },
	},
	{
		Name:   "commuting",
		Prompt: commutingPrompt,
		Criteria: func(t *Trip) []any {
			return []any{t.Destination, t.DateFrom.Format("2006-01-02"), t.DateTo.Format("2006-01-02"), t.Budget, t.Requirements}
		},
		Condition: func(t *Trip) bool { return t.Result.Commuting == "" },
	},
}

// FieldConfig represents a configuration required to generate a single field of a trip result.
type FieldConfig struct {
	Name      string
	Prompt    string
	Criteria  func(*Trip) []any
	Condition func(*Trip) bool
}

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
	errs := make(chan error, len(requestResultConf))

	// Process each field concurrently
	for _, c := range requestResultConf {
		if c.Condition(trip) {
			wg.Add(1)
			go func(c FieldConfig) {
				defer wg.Done()
				slog.Info("Generating field for trip", "trip_id", trip.ID, "field", c.Name)

				field, err := generateField(aiClient, c.Prompt, c.Criteria(trip)...)
				if err != nil {
					errs <- fmt.Errorf("generate %s: %w", c.Name, err)
					return
				}

				if err = tripsRepository.UpdateResult(
					trip.Result.ID.String(),
					map[string]any{c.Name: field},
				); err != nil {
					errs <- fmt.Errorf("update result %s: %w", c.Name, err)
					return
				}
			}(c)
		}
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
		slog.Error("Failed to generate result for trip", "trip_id", trip.ID, "errors", cerr)
		err = tripsRepository.Update(
			tripID,
			map[string]any{"status": FailedStatus.String()},
		)
		if err != nil {
			return fmt.Errorf("update trip status: %w", err)
		}
		return cerr
	}

	slog.Info("Generated result for trip", "trip_id", trip.ID)
	err = tripsRepository.Update(
		tripID,
		map[string]any{"status": CompletedStatus.String()},
	)
	if err != nil {
		return fmt.Errorf("update trip status: %w", err)
	}

	return nil
}

// generateField creates a prompt for a given field, sends it to the AI client and returns the response.
func generateField(aiClient ai.AIClient, prompt string, criteria ...any) (string, error) {
	prompt = fmt.Sprintf(basePrompt+prompt, criteria...)
	return aiClient.Send(prompt)
}

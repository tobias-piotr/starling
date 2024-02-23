package trips

import (
	"fmt"

	"starling/internal/ai"
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
	// TODO: Fetch it with the result
	trip, err := tripsRepository.Get(tripID)
	if err != nil {
		return err
	}

	// TODO: Check only empty fields and process those

	summary, err := generateSummary(aiClient, trip)
	if err != nil {
		return err
	}

	// TODO: Save the field into the result
	fmt.Println(summary)

	// TODO: Set the status to completed/failed

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

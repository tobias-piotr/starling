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
`
)

func GenerateSummary(client ai.AIClient, trip *Trip) (string, error) {
	prompt := fmt.Sprintf(
		basePrompt+summaryPrompt,
		trip.Origin,
		trip.Destination,
		trip.DateFrom.Format("2006-01-02"),
		trip.DateTo.Format("2006-01-02"),
		trip.Budget,
		trip.Requirements,
	)
	return client.Send(prompt)
}

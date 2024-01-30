package trips

type AIClient interface {
	Send(prompt string) (string, error)
}

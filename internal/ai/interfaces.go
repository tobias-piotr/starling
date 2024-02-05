package ai

type AIClient interface {
	Send(prompt string) (string, error)
}

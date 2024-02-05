package ai

type OpenAIClient struct{}

func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{}
}

func (c *OpenAIClient) Send(prompt string) (string, error) {
	return "Response from the AI", nil
}

package ai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(key string) *OpenAIClient {
	return &OpenAIClient{client: openai.NewClient(key)}
}

func (c *OpenAIClient) Send(prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("send prompt: %w", err)
	}
	return resp.Choices[0].Message.Content, nil
}

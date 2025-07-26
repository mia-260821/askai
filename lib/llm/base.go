package llm

import (
	"context"
	"fmt"
)

type Provider = string

const (
	ProviderGemini Provider = "gemini"
	ProviderOpenai Provider = "openai"
)

type Client interface {
	GetChatClient() ChatClient
}

type Session interface {
	Send(text string) (<-chan string, error)
}

type ChatClient interface {
	NewSession(ctx context.Context) (Session, error)
}

func NewClient(provider Provider, model string, apiKey string) (Client, error) {
	switch provider {
	case ProviderOpenai:
		panic("openai provider not implemented yet")
	case ProviderGemini:
		return NewGeminiClient(apiKey, model)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

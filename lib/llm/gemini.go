package llm

import (
	"context"
	"fmt"
	"google.golang.org/genai"
)

type (
	GeminiClient struct {
		client *genai.Client
		model  string
	}

	GeminiSession struct {
		chat *genai.Chat
		ctx  context.Context
	}
)

func (s *GeminiSession) Send(text string) (<-chan string, error) {
	ch := make(chan string)
	go func() {
		defer close(ch)

		resSeq := s.chat.SendStream(s.ctx, genai.NewPartFromText(text))
		resSeq(func(res *genai.GenerateContentResponse, err error) bool {
			if err != nil {
				return false
			}
			ch <- res.Text()
			return true
		})
	}()
	return ch, nil
}

func (c *GeminiClient) GetChatClient() ChatClient {
	return c
}

func (c *GeminiClient) NewSession(ctx context.Context) (Session, error) {
	cfg := genai.GenerateContentConfig{}
	history := make([]*genai.Content, 0)
	chat, err := c.client.Chats.Create(ctx, c.model, &cfg, history)
	if err != nil {
		return nil, fmt.Errorf("failed to start a chat %s", err.Error())
	}
	return &GeminiSession{chat: chat, ctx: ctx}, nil
}

func NewGeminiClient(apiKey string, model string) (Client, error) {
	ctx := context.Background()
	cfg := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}
	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

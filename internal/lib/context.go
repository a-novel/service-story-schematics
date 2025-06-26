package lib

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/a-novel/service-story-schematics/config"
)

func NewAgoraContext(parentCTX context.Context, dsn string) (context.Context, error) {
	ctx, err := NewPostgresContext(parentCTX, dsn)
	if err != nil {
		return parentCTX, fmt.Errorf("create postgres context: %w", err)
	}

	return NewOpenaiContext(ctx), nil
}

type OpenaiCtxKey struct{}

func NewOpenaiContext(parentCTX context.Context) context.Context {
	client := openai.NewClient(
		option.WithAPIKey(config.Groq.APIKey),
		option.WithBaseURL(config.Groq.BaseURL),
	)

	return context.WithValue(parentCTX, OpenaiCtxKey{}, &client)
}

func OpenAIClient(ctx context.Context) *openai.Client {
	client, ok := ctx.Value(OpenaiCtxKey{}).(*openai.Client)
	if !ok {
		return nil
	}

	return client
}

package lib

import (
	"context"
	"fmt"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/service-story-schematics/migrations"
)

func NewAgoraContext(parentCTX context.Context) (context.Context, error) {
	ctx, err := pgctx.NewContext(parentCTX, &migrations.Migrations)
	if err != nil {
		return nil, fmt.Errorf("create postgres context: %w", err)
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

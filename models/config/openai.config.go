package config

import (
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/samber/lo"
)

type OpenAI struct {
	Model   openai.ChatModel `json:"model"   yaml:"model"`
	APIKey  string           `json:"apiKey"  yaml:"apiKey"`
	BaseURL string           `json:"baseURL" yaml:"baseURL"`
}

func (preset OpenAI) Client() *openai.Client {
	return lo.ToPtr(openai.NewClient(
		option.WithAPIKey(preset.APIKey),
		option.WithBaseURL(preset.BaseURL),
	))
}

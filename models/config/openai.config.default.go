package config

import (
	"github.com/openai/openai-go/v2"

	"github.com/a-novel/golib/config"
)

const (
	OpenAIModel   openai.ChatModel = "meta-llama/llama-4-maverick-17b-128e-instruct"
	OpenAIBaseURL openai.ChatModel = "https://api.groq.com/openai/v1"
)

var OpenAIPresetDefault = OpenAI{
	APIKey:  getEnv("OPENAI_TOKEN"),
	Model:   config.LoadEnv(getEnv("OPENAI_MODEL"), OpenAIModel, config.StringParser),
	BaseURL: config.LoadEnv(getEnv("OPENAI_BASE_URL"), OpenAIBaseURL, config.StringParser),
}

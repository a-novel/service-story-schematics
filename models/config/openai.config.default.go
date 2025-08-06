package config

import (
	"github.com/openai/openai-go"

	"github.com/a-novel/golib/config"
)

const (
	OpenAIModel   openai.ChatModel = "openai/gpt-oss-120b"
	OpenAIBaseURL openai.ChatModel = "https://api.groq.com/openai/v1"
)

var OpenAIPresetDefault = OpenAI{
	APIKey:  getEnv("OPENAI_TOKEN"),
	Model:   config.LoadEnv(getEnv("OPENAI_MODEL"), OpenAIModel, config.StringParser),
	BaseURL: config.LoadEnv(getEnv("OPENAI_BASE_URL"), OpenAIBaseURL, config.StringParser),
}

package config

import (
	"os"

	"github.com/openai/openai-go"

	"github.com/a-novel/golib/config"
)

const (
	OpenAIModel   openai.ChatModel = "meta-llama/llama-4-maverick-17b-128e-instruct"
	OpenAIBaseURL openai.ChatModel = "https://api.groq.com/openai/v1"
)

var OpenAIPresetDefault = OpenAI{
	APIKey:  os.Getenv("OPENAI_TOKEN"),
	Model:   config.LoadEnv(os.Getenv("OPENAI_MODEL"), OpenAIModel, config.StringParser),
	BaseURL: config.LoadEnv(os.Getenv("OPENAI_BASE_URL"), OpenAIBaseURL, config.StringParser),
}

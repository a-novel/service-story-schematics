package config

import (
	_ "embed"

	"github.com/openai/openai-go"

	"github.com/a-novel-kit/configurator"
)

//go:embed groq.yaml
var groqFile []byte

type GroqType struct {
	APIKey  string           `yaml:"apiKey"`
	Model   openai.ChatModel `yaml:"model"`
	BaseURL string           `yaml:"baseURL"`
}

var Groq = configurator.NewLoader[GroqType](Loader).MustLoad(
	configurator.NewConfig("", groqFile),
)

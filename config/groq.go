package config

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"
	groqmodels "github.com/a-novel-kit/golm/bindings/groq/models"
)

//go:embed groq.yaml
var groqFile []byte

type GroqType struct {
	APIKey string           `yaml:"apiKey"`
	Model  groqmodels.Model `yaml:"model"`
}

var Groq = configurator.NewLoader[GroqType](Loader).MustLoad(
	configurator.NewConfig("", groqFile),
)

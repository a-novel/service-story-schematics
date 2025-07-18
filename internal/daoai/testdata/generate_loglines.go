package testdata

import (
	_ "embed"
	"github.com/a-novel/golib/config"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/goccy/go-yaml"
)

//go:embed generate_loglines.en.yaml
var generateLoglinesEnFile []byte

//go:embed generate_loglines.fr.yaml
var generateLoglinesFrFile []byte

type GenerateLoglinesTestCase struct {
	Count int    `yaml:"count"`
	Theme string `yaml:"theme"`
}

type GenerateLoglinesPromptsType struct {
	Cases      map[string]GenerateLoglinesTestCase `yaml:"cases"`
	CheckAgent struct {
		Themed string `yaml:"themed"`
		Random string `yaml:"random"`
	} `yaml:"checkAgent"`
}

var GenerateLoglinesPromptEN = config.MustUnmarshal[GenerateLoglinesPromptsType](yaml.Unmarshal, generateLoglinesEnFile)

var GenerateLoglinesPromptFR = config.MustUnmarshal[GenerateLoglinesPromptsType](yaml.Unmarshal, generateLoglinesFrFile)

var GenerateLoglinesPrompts = map[models.Lang]GenerateLoglinesPromptsType{
	models.LangEN: GenerateLoglinesPromptEN,
	models.LangFR: GenerateLoglinesPromptFR,
}

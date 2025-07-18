package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed generate_loglines.en.yaml
var generateLoglinesEnFile []byte

//go:embed generate_loglines.fr.yaml
var generateLoglinesFrFile []byte

type GenerateLoglinessType struct {
	System struct {
		Themed string `yaml:"themed"`
		Random string `yaml:"random"`
	} `yaml:"system"`
}

var GenerateLoglinesEN = config.MustUnmarshal[GenerateLoglinessType](yaml.Unmarshal, generateLoglinesEnFile)

var GenerateLoglinesFR = config.MustUnmarshal[GenerateLoglinessType](yaml.Unmarshal, generateLoglinesFrFile)

var GenerateLoglines = map[models.Lang]GenerateLoglinessType{
	models.LangEN: GenerateLoglinesEN,
	models.LangFR: GenerateLoglinesFR,
}

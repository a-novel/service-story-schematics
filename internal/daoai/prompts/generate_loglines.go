package prompts

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
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

var GenerateLoglinesEN = configurator.NewLoader[GenerateLoglinessType](config.Loader).MustLoad(
	configurator.NewConfig("", generateLoglinesEnFile),
)

var GenerateLoglinesFR = configurator.NewLoader[GenerateLoglinessType](config.Loader).MustLoad(
	configurator.NewConfig("", generateLoglinesFrFile),
)

var GenerateLoglines = map[models.Lang]GenerateLoglinessType{
	models.LangEN: GenerateLoglinesEN,
	models.LangFR: GenerateLoglinesFR,
}

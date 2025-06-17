package testdata

import (
	_ "embed"
	"github.com/a-novel-kit/configurator"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed utils.en.yaml
var utilsEnFile []byte

//go:embed utils.fr.yaml
var utilsFrFile []byte

type UtilsPromptsType struct {
	CheckAgent struct {
		System string `yaml:"system"`
		Expect string `yaml:"expect"`
	} `yaml:"checkAgent"`
}

var UtilsPromptEN = configurator.NewLoader[UtilsPromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", utilsEnFile),
)

var UtilsPromptFR = configurator.NewLoader[UtilsPromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", utilsFrFile),
)

var UtilsPrompts = map[models.Lang]UtilsPromptsType{
	models.LangEN: UtilsPromptEN,
	models.LangFR: UtilsPromptFR,
}

package testdata

import (
	_ "embed"
	"github.com/a-novel-kit/configurator"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed expand_logline.en.yaml
var expandLoglineEnFile []byte

//go:embed expand_logline.fr.yaml
var expandLoglineFrFile []byte

type ExpandLoglineTestCase struct {
	Logline string `yaml:"logline"`
}

type ExpandLoglinePromptsType struct {
	Cases      map[string]ExpandLoglineTestCase `yaml:"cases"`
	CheckAgent string                           `yaml:"checkAgent"`
}

var ExpandLoglinePromptEN = configurator.NewLoader[ExpandLoglinePromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", expandLoglineEnFile),
)

var ExpandLoglinePromptFR = configurator.NewLoader[ExpandLoglinePromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", expandLoglineFrFile),
)

var ExpandLoglinePrompts = map[models.Lang]ExpandLoglinePromptsType{
	models.LangEN: ExpandLoglinePromptEN,
	models.LangFR: ExpandLoglinePromptFR,
}

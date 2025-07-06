package prompts

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

type ExpandLoglinesType struct {
	System string `yaml:"system"`
}

var ExpandLoglineEN = configurator.NewLoader[ExpandLoglinesType](config.Loader).MustLoad(
	configurator.NewConfig("", expandLoglineEnFile),
)

var ExpandLoglineFR = configurator.NewLoader[ExpandLoglinesType](config.Loader).MustLoad(
	configurator.NewConfig("", expandLoglineFrFile),
)

var ExpandLogline = map[models.Lang]ExpandLoglinesType{
	models.LangEN: ExpandLoglineEN,
	models.LangFR: ExpandLoglineFR,
}

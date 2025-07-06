package prompts

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed expand_beat.en.yaml
var expandBeatEnFile []byte

//go:embed expand_beat.fr.yaml
var expandBeatFrFile []byte

type ExpandBeatsType struct {
	System string `yaml:"system"`
	Input1 string `yaml:"input1"`
	Input2 string `yaml:"input2"`
}

var ExpandBeatEN = configurator.NewLoader[ExpandBeatsType](config.Loader).MustLoad(
	configurator.NewConfig("", expandBeatEnFile),
)

var ExpandBeatFR = configurator.NewLoader[ExpandBeatsType](config.Loader).MustLoad(
	configurator.NewConfig("", expandBeatFrFile),
)

var ExpandBeat = map[models.Lang]ExpandBeatsType{
	models.LangEN: ExpandBeatEN,
	models.LangFR: ExpandBeatFR,
}

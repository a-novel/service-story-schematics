package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

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

var ExpandBeatEN = config.MustUnmarshal[ExpandBeatsType](yaml.Unmarshal, expandBeatEnFile)

var ExpandBeatFR = config.MustUnmarshal[ExpandBeatsType](yaml.Unmarshal, expandBeatFrFile)

var ExpandBeat = map[models.Lang]ExpandBeatsType{
	models.LangEN: ExpandBeatEN,
	models.LangFR: ExpandBeatFR,
}

package testdata

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

type ExpandBeatTestCase struct {
	Logline   string           `yaml:"logline"`
	Plan      models.StoryPlan `yaml:"plan"`
	Beats     []models.Beat    `yaml:"beats"`
	TargetKey string           `yaml:"targetKey"`
}

type ExpandBeatPromptsType struct {
	Cases      map[string]ExpandBeatTestCase `yaml:"cases"`
	CheckAgent string                        `yaml:"checkAgent"`
}

var ExpandBeatPromptEN = configurator.NewLoader[ExpandBeatPromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", expandBeatEnFile),
)

var ExpandBeatPromptFR = configurator.NewLoader[ExpandBeatPromptsType](config.Loader).MustLoad(
	configurator.NewConfig("", expandBeatFrFile),
)

var ExpandBeatPrompts = map[models.Lang]ExpandBeatPromptsType{
	models.LangEN: ExpandBeatPromptEN,
	models.LangFR: ExpandBeatPromptFR,
}

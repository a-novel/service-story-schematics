package prompts

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
)

//go:embed en.yaml
var en []byte

type Prompts struct {
	GenerateLogline struct {
		System struct {
			Themed string `yaml:"themed"`
			Random string `yaml:"random"`
		} `yaml:"system"`
	} `yaml:"generateLoglines"`
	ExpandLogline string `yaml:"expandLogline"`

	StoryPlan string `yaml:"storyPlan"`

	GenerateBeatsSheet string `yaml:"generateBeatsSheet"`
	ExpandBeat         struct {
		System string `yaml:"system"`
		Input1 string `yaml:"input1"`
		Input2 string `yaml:"input2"`
	} `yaml:"expandBeat"`
	RegenerateBeats struct {
		System string `yaml:"system"`
		Input1 string `yaml:"input1"`
		Input2 string `yaml:"input2"`
	} `yaml:"regenerateBeats"`
}

type TranslatedPrompts struct {
	En Prompts `yaml:"en"`
}

var Config = TranslatedPrompts{
	En: configurator.NewLoader[Prompts](config.Loader).MustLoad(configurator.NewConfig("", en)),
}

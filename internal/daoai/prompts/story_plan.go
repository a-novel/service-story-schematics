package prompts

import (
	_ "embed"

	"github.com/a-novel-kit/configurator"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/models"
)

//go:embed story_plan.en.yaml
var storyPlanEnFile []byte

//go:embed story_plan.fr.yaml
var storyPlanFrFile []byte

type StoryPlansType struct {
	System string `yaml:"system"`
}

var StoryPlanEN = configurator.NewLoader[StoryPlansType](config.Loader).MustLoad(
	configurator.NewConfig("", storyPlanEnFile),
)

var StoryPlanFR = configurator.NewLoader[StoryPlansType](config.Loader).MustLoad(
	configurator.NewConfig("", storyPlanFrFile),
)

var StoryPlan = map[models.Lang]StoryPlansType{
	models.LangEN: StoryPlanEN,
	models.LangFR: StoryPlanFR,
}

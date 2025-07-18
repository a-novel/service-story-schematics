package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed story_plan.en.yaml
var storyPlanEnFile []byte

//go:embed story_plan.fr.yaml
var storyPlanFrFile []byte

type StoryPlansType struct {
	System string `yaml:"system"`
}

var StoryPlanEN = config.MustUnmarshal[StoryPlansType](yaml.Unmarshal, storyPlanEnFile)

var StoryPlanFR = config.MustUnmarshal[StoryPlansType](yaml.Unmarshal, storyPlanFrFile)

var StoryPlan = map[models.Lang]StoryPlansType{
	models.LangEN: StoryPlanEN,
	models.LangFR: StoryPlanFR,
}

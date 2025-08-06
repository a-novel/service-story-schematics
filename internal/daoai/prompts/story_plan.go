package prompts

import (
	_ "embed"

	"github.com/goccy/go-yaml"

	"github.com/a-novel/golib/config"
)

//go:embed story_plan.en.yaml
var storyPlanEnFile []byte

type StoryPlansType struct {
	System string `yaml:"system"`
}

var StoryPlan = config.MustUnmarshal[StoryPlansType](yaml.Unmarshal, storyPlanEnFile)

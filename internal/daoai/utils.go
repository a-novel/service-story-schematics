package daoai

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

var StoryPlanPrompt = template.Must(template.New("EN").Parse(prompts.Config.En.StoryPlan))

func StoryPlanToPrompt(tName string, data models.StoryPlan) (string, error) {
	var sb strings.Builder

	err := StoryPlanPrompt.ExecuteTemplate(&sb, tName, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute story plan prompt template: %w", err)
	}

	return sb.String(), nil
}

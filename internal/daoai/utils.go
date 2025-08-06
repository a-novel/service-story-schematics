package daoai

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

var StoryPlanPrompt = template.Must(template.New("").Parse(prompts.StoryPlan.System))

func StoryPlanToPrompt(data models.StoryPlan) (string, error) {
	var sb strings.Builder

	err := StoryPlanPrompt.ExecuteTemplate(&sb, "", data)
	if err != nil {
		return "", fmt.Errorf("failed to execute story plan prompt template: %w", err)
	}

	return sb.String(), nil
}

func ForceNextAnswerLocale(locale models.Lang, system string) string {
	translatePrompt, ok := prompts.Langs[locale]
	if !ok {
		return system
	}

	return strings.Join([]string{system, translatePrompt}, "\n")
}

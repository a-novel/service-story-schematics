package daoai

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

var StoryPlanPrompt = RegisterTemplateLocales(prompts.StoryPlan[models.LangEN].System, map[models.Lang]string{
	models.LangEN: prompts.StoryPlan[models.LangEN].System,
	models.LangFR: prompts.StoryPlan[models.LangFR].System,
})

func StoryPlanToPrompt(lang models.Lang, data models.StoryPlan) (string, error) {
	var sb strings.Builder

	err := StoryPlanPrompt.ExecuteTemplate(&sb, lang.String(), data)
	if err != nil {
		return "", fmt.Errorf("failed to execute story plan prompt template: %w", err)
	}

	return sb.String(), nil
}

func RegisterTemplateLocales(def string, data map[models.Lang]string) *template.Template {
	t := template.Must(template.New("").Parse(def))

	for lang, prompt := range data {
		template.Must(t.New(lang.String()).Parse(prompt))
	}

	return t
}

package daoai

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

var StoryPlanPrompt = RegisterTemplateLocales(prompts.Config.En.StoryPlan, map[models.Lang]string{
	models.LangEN: prompts.Config.En.StoryPlan,
	models.LangFR: prompts.Config.Fr.StoryPlan,
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

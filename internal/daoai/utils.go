package daoai

import (
	"strings"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

func ForceNextAnswerLocale(locale models.Lang, system string) string {
	translatePrompt, ok := prompts.Langs[locale]
	if !ok {
		return system
	}

	return strings.Join([]string{system, translatePrompt}, "\n")
}

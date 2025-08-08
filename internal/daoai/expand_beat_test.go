package daoai_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

func TestExpandBeat(t *testing.T) {
	const errorMsg = "The new beat does not expand on the original beat.\n\n" +
		"new sheet:\n\n%s\n\noriginal beat:\n\n%s"

	repository := daoai.NewExpandBeatRepository(&config.OpenAIPresetDefault)

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.ExpandBeatPrompt
			plan := storyplanmodel.SaveTheCat[lang].Pick("openingImage", "themeStated", "setup", "catalyst", "debate")

			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					resp, err := repository.ExpandBeat(t.Context(), daoai.ExpandBeatRequest{
						Logline:   testCase.Logline,
						Beats:     testCase.Beats,
						Plan:      plan,
						Lang:      lang,
						TargetKey: testCase.TargetKey,
						UserID:    TestUser,
					})
					require.NoError(t, err)

					require.NotNil(t, resp)

					original, ok := lo.Find(testCase.Beats, func(item models.Beat) bool {
						return item.Key == testCase.TargetKey
					})
					require.True(t, ok)

					require.NotEqual(t, original.Content, resp.Content)

					CheckAgent(
						t,
						fmt.Sprintf(data.CheckAgent, resp.Content, original),
						fmt.Sprintf(errorMsg, resp.Content, original),
					)
					CheckLang(t, lang, resp.Content)
				})
			}
		})
	}
}

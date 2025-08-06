package daoai_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

func TestGenerateLoglines(t *testing.T) {
	const errorMsgThemed = "The greater AI decreted that this logline:\n\n%s\n\nDoes not match this theme:\n\n%s"

	const errorMsgRandom = "The greater AI decreted that this is not a valid logline for a story:\n\n%s"

	repository := daoai.NewGenerateLoglinesRepository(&config.OpenAIPresetDefault)

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.GenerateLoglinesPrompt

			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					loglines, err := repository.GenerateLoglines(t.Context(), daoai.GenerateLoglinesRequest{
						Count:  testCase.Count,
						Theme:  testCase.Theme,
						UserID: TestUser,
						Lang:   lang,
					})
					require.NoError(t, err)

					require.NotNil(t, loglines)
					require.Len(t, loglines, testCase.Count)

					for _, logline := range loglines {
						require.NotEmpty(t, logline.Name)
						require.NotEmpty(t, logline.Content)

						if testCase.Theme != "" {
							CheckAgent(
								t,
								fmt.Sprintf(data.CheckAgent.Themed, logline.Content, testCase.Theme),
								fmt.Sprintf(errorMsgThemed, logline.Content, testCase.Theme),
							)
						} else {
							CheckAgent(
								t,
								fmt.Sprintf(data.CheckAgent.Random, logline.Content),
								fmt.Sprintf(errorMsgRandom, logline.Content),
							)
						}

						CheckLang(t, lang, logline.Content)
					}
				})
			}
		})
	}
}

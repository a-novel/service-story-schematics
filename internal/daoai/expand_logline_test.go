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

func TestExpandLogline(t *testing.T) {
	const errorMsg = "The greater AI decreted that this logline:\n\n%s\n\nDoes not expand this one:\n\n%s"

	repository := daoai.NewExpandLoglineRepository(&config.OpenAIPresetDefault)

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.ExpandLoglinePrompts[lang]

			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					resp, err := repository.ExpandLogline(t.Context(), daoai.ExpandLoglineRequest{
						Logline: testCase.Logline,
						Lang:    lang,
						UserID:  TestUser,
					})
					require.NoError(t, err)

					require.NotNil(t, resp)

					require.NotEmpty(t, resp.Name)
					require.NotEmpty(t, resp.Content)

					CheckAgent(
						t, lang,
						fmt.Sprintf(data.CheckAgent, resp.Content, testCase.Logline),
						fmt.Sprintf(errorMsg, resp.Content, testCase.Logline),
					)
				})
			}
		})
	}
}

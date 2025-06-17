package daoai_test

import (
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
)

func TestExpandLogline(t *testing.T) {
	const errorMsg = "The greater AI decreted that this logline:\n\n%s\n\nDoes not expand this one:\n\n%s"

	repository := daoai.NewExpandLoglineRepository()

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.ExpandLoglinePrompts[lang]

			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					ctx := lib.NewOpenaiContext(t.Context())

					resp, err := repository.ExpandLogline(ctx, daoai.ExpandLoglineRequest{
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

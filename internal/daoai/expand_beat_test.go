package daoai_test

import (
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

func TestExpandBeat(t *testing.T) {
	const errorMsg = "The new beat does not expand on the original beat.\n\n" +
		"new sheet:\n\n%s\n\noriginal beat:\n\n%s"

	repository := daoai.NewExpandBeatRepository()

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.ExpandBeatPrompts[lang]
			
			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					ctx := lib.NewOpenaiContext(t.Context())

					resp, err := repository.ExpandBeat(ctx, daoai.ExpandBeatRequest{
						Logline:   testCase.Logline,
						Beats:     testCase.Beats,
						Plan:      testCase.Plan,
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
						t, lang,
						fmt.Sprintf(data.CheckAgent, resp.Content, original),
						fmt.Sprintf(errorMsg, resp.Content, original),
					)
				})
			}
		})
	}
}

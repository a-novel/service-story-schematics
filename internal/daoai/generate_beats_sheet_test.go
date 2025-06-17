package daoai_test

import (
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

func TestGenerateBeatsSheet(t *testing.T) {
	const errorMsg = "The below beats sheet does not form a coherent story about the below logline.\n\n" +
		"beats sheet:\n\n%s\n\nlogline:\n\n%s"

	repository := daoai.NewGenerateBeatsSheetRepository()

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.GenerateBeatsSheetPrompts[lang]

			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					ctx := lib.NewOpenaiContext(t.Context())

					beatsSheet, err := repository.GenerateBeatsSheet(ctx, daoai.GenerateBeatsSheetRequest{
						Logline: testCase.Logline,
						Plan:    testCase.Plan,
						UserID:  TestUser,
						Lang:    lang,
					})
					require.NoError(t, err)

					require.NotNil(t, beatsSheet)
					require.Len(t, beatsSheet, len(testCase.Plan.Beats))

					aggregated := strings.Join(lo.Map(beatsSheet, func(item models.Beat, _ int) string {
						return item.Title + "\n" + item.Content
					}), "\n\n")

					CheckAgent(
						t, lang,
						fmt.Sprintf(data.CheckAgent, aggregated, testCase.Logline),
						fmt.Sprintf(errorMsg, aggregated, testCase.Logline),
					)
				})
			}
		})
	}
}

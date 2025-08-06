package daoai_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

func TestRegenerateBeats(t *testing.T) {
	const errorMsg = "The below beats sheet does not form a coherent story about the below logline.\n\n" +
		"beats sheet:\n\n%s\n\nlogline:\n\n%s"

	repository := daoai.NewRegenerateBeatsRepository(&config.OpenAIPresetDefault)

	for _, lang := range []models.Lang{models.LangEN, models.LangFR} {
		t.Run(lang.String(), func(t *testing.T) {
			t.Parallel()

			data := testdata.RegenerateBeatsPrompt

			for name, testCase := range data.Cases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					beatsSheet, err := repository.RegenerateBeats(t.Context(), daoai.RegenerateBeatsRequest{
						Logline:        testCase.Logline,
						Beats:          testCase.Beats,
						Plan:           testCase.Plan,
						RegenerateKeys: testCase.RegenerateKeys,
						UserID:         TestUser,
						Lang:           lang,
					})
					require.NoError(t, err)

					require.NotNil(t, beatsSheet)
					require.Len(t, beatsSheet, len(testCase.Plan.Beats))

					var aggregatedBeats, aggregatedNewBeats []string

					for i, beat := range beatsSheet {
						inlinedBeat := beat.Title + "\n" + beat.Content
						aggregatedBeats = append(aggregatedBeats, inlinedBeat)

						if lo.Contains(testCase.RegenerateKeys, beat.Key) {
							require.NotEqual(t, beat.Content, testCase.Beats[i].Content)

							aggregatedNewBeats = append(aggregatedNewBeats, inlinedBeat)
						} else {
							require.Equal(t, beat.Content, testCase.Beats[i].Content)
						}
					}

					CheckAgent(
						t,
						fmt.Sprintf(data.CheckAgent, strings.Join(aggregatedBeats, "\n"), testCase.Logline),
						fmt.Sprintf(errorMsg, strings.Join(aggregatedBeats, "\n"), testCase.Logline),
					)
					CheckLang(t, lang, strings.Join(aggregatedNewBeats, "\n"))
				})
			}
		})
	}
}

package daoai_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golm"
	groqbinding "github.com/a-novel-kit/golm/bindings/groq"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/daoai"
)

func TestGenerateLoglines(t *testing.T) {
	testCases := []struct {
		name string

		request daoai.GenerateLoglinesRequest
	}{
		{
			name: "Success",

			request: daoai.GenerateLoglinesRequest{
				Count:  3,
				Theme:  "Sci-Fi",
				UserID: TestUser,
			},
		},
		{
			name: "Success/AnotherTheme",

			request: daoai.GenerateLoglinesRequest{
				Count:  3,
				Theme:  "old school detective story",
				UserID: TestUser,
			},
		},
		{
			name: "Random",

			request: daoai.GenerateLoglinesRequest{
				Count:  3,
				UserID: TestUser,
			},
		},
	}

	binding := groqbinding.New(config.Groq.APIKey, config.Groq.Model)

	repository := daoai.NewGenerateLoglinesRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := golm.WithContext(t.Context(), binding)

			loglines, err := repository.GenerateLoglines(ctx, testCase.request)
			require.NoError(t, err)

			require.NotNil(t, loglines)
			require.Len(t, loglines, testCase.request.Count)

			for _, logline := range loglines {
				require.NotEmpty(t, logline.Name)
				require.NotEmpty(t, logline.Content)

				if testCase.request.Theme != "" {
					CheckAgent(
						t,
						fmt.Sprintf(
							"Does this logline:\n\n%s\n\nMatches this theme:\n\n%s",
							logline.Content, testCase.request.Theme,
						),
						fmt.Sprintf(
							"The greater AI decreted that this logline:\n\n%s\n\nDoes not match this theme:\n\n%s",
							logline.Content, testCase.request.Theme,
						),
					)
				} else {
					CheckAgent(
						t,
						"Can this be used as a logline for a story ?\n\n"+logline.Content,
						"The greater AI decreted that this is not a valid logline for a story:\n\n"+logline.Content,
					)
				}
			}
		})
	}
}

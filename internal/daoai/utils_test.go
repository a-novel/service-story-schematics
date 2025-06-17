package daoai_test

import (
	"github.com/a-novel/service-story-schematics/internal/daoai/testdata"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

const TestUser = "agora-story-schematics-test"

// Remove punctuation and other extra characters from the check agent response.
var checkAgentFormatter = regexp.MustCompile("[^a-zA-Z0-9 ]")

// CheckAgent validates the response of an OpenAI agent to a given prompt.
func CheckAgent(t *testing.T, lang models.Lang, prompt, message string) {
	t.Helper()

	ctx := lib.NewOpenaiContext(t.Context())
	chatCompletion, err := lib.OpenAIClient(ctx).Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       config.Groq.Model,
		Temperature: param.NewOpt(0.2),
		User:        param.NewOpt(TestUser),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(testdata.UtilsPrompts[lang].CheckAgent.System),
			openai.UserMessage(prompt),
		},
	})

	require.NoError(t, err)

	// Prevent flakiness issues because LLMS just dont give a fuck.
	formattedResponse := chatCompletion.Choices[0].Message.Content
	formattedResponse = strings.TrimSpace(formattedResponse)
	formattedResponse = strings.ToUpper(formattedResponse)
	formattedResponse = checkAgentFormatter.ReplaceAllString(formattedResponse, "")
	require.Equal(t, testdata.UtilsPrompts[lang].CheckAgent.Expect, formattedResponse, message)
}

func TestStoryPlanToPrompt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		template models.Lang
		plan     models.StoryPlan

		expect string
	}{
		{
			name: "En",

			template: "en",

			plan: models.StoryPlan{
				Description: "Test Description.",
				Beats: []models.BeatDefinition{
					{
						Name:      "Beat 1",
						MinScenes: 1,
						Key:       "beat1",
						KeyPoints: []string{
							"Beat 1 - Key Point 1",
							"Beat 1 - Key Point 2",
						},
						Purpose: "Beat 1 - Purpose",
					},
					{
						Name:      "Beat 2",
						MinScenes: 2,
						MaxScenes: 5,
						Key:       "beat2",
						KeyPoints: []string{
							"Beat 2 - Key Point 1",
						},
						Purpose: "Beat 2 - Purpose",
					},
				},
			},

			expect: `Test Description.

Here's a detailed breakdown with minimum scenes and key points for each beat:

Beat 1 (1 scene)
JSON Key: beat1
Key points:
  - Beat 1 - Key Point 1
  - Beat 1 - Key Point 2
Purpose: Beat 1 - Purpose

Beat 2 (2 to 5 scenes)
JSON Key: beat2
Key points:
  - Beat 2 - Key Point 1
Purpose: Beat 2 - Purpose

This concludes the breakdown. Below are important things for you to consider.

Focus on Essentials:
Ensure each scene serves a clear purpose and advances the plot.

Avoid Redundancy:
Eliminate unnecessary scenes that don't contribute to character development or plot progression.

Balance Pacing:
Allocate scenes strategically to maintain engagement throughout the story.

Character Development:
Ensure each scene contributes to character growth and progression.`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			prompt, err := daoai.StoryPlanToPrompt(testCase.template, testCase.plan)
			require.NoError(t, err)
			require.Equal(t, strings.TrimSpace(testCase.expect), strings.TrimSpace(prompt))
		})
	}
}

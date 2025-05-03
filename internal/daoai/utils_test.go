package daoai_test

import (
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golm"
	groqbinding "github.com/a-novel-kit/golm/bindings/groq"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

const TestUser = "agora-story-schematics-test"

func CheckAgent(t *testing.T, prompt, message string) {
	t.Helper()

	binding := groqbinding.New(config.Groq.APIKey, config.Groq.Model)
	chat := golm.NewChat(binding)
	chat.SetSystem(golm.NewSystemMessage(`You are a tester. Just answer "YES" or "NO", nothing else.`))

	requestMessage := golm.NewUserMessage(prompt)

	params := golm.CompletionParams{
		Temperature: lo.ToPtr(0.2),
	}

	resp, err := chat.Completion(t.Context(), requestMessage, params)
	require.NoError(t, err)

	require.Equal(t, "YES", resp.GetContent(), message)
}

func TestStoryPlanToPrompt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		template string
		plan     models.StoryPlan

		expect string
	}{
		{
			name: "En",

			template: "EN",

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

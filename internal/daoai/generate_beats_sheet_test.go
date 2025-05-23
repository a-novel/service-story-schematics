package daoai_test

import (
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

func TestGenerateBeatsSheet(t *testing.T) {
	testCases := []struct {
		name string

		request daoai.GenerateBeatsSheetRequest
	}{
		{
			name: "Success",

			request: daoai.GenerateBeatsSheetRequest{
				Logline: `The Aurora Initiative

As a team of scientists discover a way to harness the energy of a nearby supernova, they must also contend with the 
implications of altering the course of human history and the emergence of a new, technologically advanced world order.`,
				Plan: models.StoryPlan{
					Name: "Save The Cat Simplified",
					Description: `The "Save The Cat" simplified story plan consists of 5 beats that serve as a
blueprint for crafting compelling stories.`,
					Beats: []models.BeatDefinition{
						{
							Name:      "Opening Image",
							Key:       "openingImage",
							KeyPoints: []string{"Establish the protagonist's world before the journey begins."},
							Purpose: "Sets the tone, mood, and stakes; offers a visual representation " +
								"of the starting point.",
							MinScenes: 1,
							MaxScenes: 1,
						},
						{
							Name:      "Theme Stated",
							Key:       "themeStated",
							KeyPoints: []string{"Introduce the story's central theme or moral."},
							Purpose:   "Often delivered through dialogue; foreshadows the protagonist transformation.",
							MinScenes: 1,
						},
						{
							Name: "Set-Up",
							Key:  "setup",
							KeyPoints: []string{
								"Introduce the main characters.",
								"Showcase the protagonist's flaws or challenges.",
								"Establish the stakes and the world they live in.",
							},
							Purpose:   "Builds empathy and grounds the audience in the story.",
							MinScenes: 3,
							MaxScenes: 5,
						},
						{
							Name:      "Catalyst",
							Key:       "catalyst",
							KeyPoints: []string{"An event that disrupts the status quo."},
							Purpose:   "Propels the protagonist into the main conflict.",
							MaxScenes: 1,
						},
						{
							Name: "Debate",
							Key:  "debate",
							KeyPoints: []string{
								"The protagonist grapples with the decision to embark on the journey.",
								"Highlights internal conflicts and fears.",
							},
							Purpose:   "Adds depth to the character and heightens tension.",
							MinScenes: 2,
							MaxScenes: 3,
						},
					},
				},
				UserID: TestUser,
			},
		},
	}

	repository := daoai.NewGenerateBeatsSheetRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := lib.NewOpenaiContext(t.Context())

			beatsSheet, err := repository.GenerateBeatsSheet(ctx, testCase.request)
			require.NoError(t, err)

			require.NotNil(t, beatsSheet)
			require.Len(t, beatsSheet, len(testCase.request.Plan.Beats))

			aggregated := strings.Join(lo.Map(beatsSheet, func(item models.Beat, _ int) string {
				return item.Title + "\n" + item.Content
			}), "\n\n")

			CheckAgent(
				t,
				fmt.Sprintf(
					"Does the below beats sheet form a coherent story about the below logline?\n\n"+
						"beats sheet:\n\n%s\n\nlogline:\n\n%s",
					aggregated, testCase.request.Logline,
				),
				fmt.Sprintf(
					"The below beats sheet does not form a coherent story about the below logline.\n\n"+
						"beats sheet:\n\n%s\n\nlogline:\n\n%s",
					aggregated, testCase.request.Logline,
				),
			)
		})
	}
}

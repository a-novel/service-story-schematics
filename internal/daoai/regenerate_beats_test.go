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

func TestRegenerateBeats(t *testing.T) {
	testCases := []struct {
		name string

		request daoai.RegenerateBeatsRequest
	}{
		{
			name: "RegenerateOne",

			request: daoai.RegenerateBeatsRequest{
				Logline: `The Aurora Initiative

As a team of scientists discover a way to harness the energy of a nearby supernova, they must also contend with the 
implications of altering the course of human history and the emergence of a new, technologically advanced world order.`,
				Plan: models.StoryPlan{
					Name: "Save The Cat Simplified",
					Description: `The "Save The Cat" simplified story plan consists of 5 beats that serve as a
blueprint for crafting compelling stories.`,
					Lang: models.LangEN,
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
				Beats: []models.Beat{
					{
						Key:   "openingImage",
						Title: "Introduction to the Scientific Community",
						Content: "A scene showcasing the team of scientists and their ordinary world, highlighting " +
							"the current limitations and struggles in the field of energy production",
					},
					{
						Key:   "themeStated",
						Title: "The Importance of Responsible Innovation",
						Content: "A conversation among the scientists discussing the ethics and potential " +
							"consequences of tapping into extraordinary energy sources, foreshadowing the themes " +
							"of the story.",
					},
					{
						Key:   "setup",
						Title: "The World on the Brink of Change",
						Content: "Introduction to the main characters, their motivations, and the world they " +
							"live in, highlighting the team's pioneering work and the world's current energy crisis.",
					},
					{
						Key:   "catalyst",
						Title: "The Discovery of the Supernova Energy Source",
						Content: "A critical experiment or finding that reveals the possibility of harnessing " +
							"the energy of a nearby supernova, initiating the team's journey into the unknown.",
					},
					{
						Key:   "debate",
						Title: "The Dilemma of Power and Responsibility",
						Content: "The team grapples with the implications and potential fallout of their " +
							"discovery, voicing concerns and conflicting views on how to proceed, and setting the " +
							"stage for character arcs.",
					},
				},
				Lang:           models.LangEN,
				RegenerateKeys: []string{"setup"},
				UserID:         TestUser,
			},
		},
		{
			name: "RegenerateMany",

			request: daoai.RegenerateBeatsRequest{
				Logline: `The Aurora Initiative

As a team of scientists discover a way to harness the energy of a nearby supernova, they must also contend with the 
implications of altering the course of human history and the emergence of a new, technologically advanced world order.`,
				Plan: models.StoryPlan{
					Name: "Save The Cat Simplified",
					Description: `The "Save The Cat" simplified story plan consists of 5 beats that serve as a
blueprint for crafting compelling stories.`,
					Lang: models.LangEN,
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
				Beats: []models.Beat{
					{
						Key:   "openingImage",
						Title: "Introduction to the Scientific Community",
						Content: "A scene showcasing the team of scientists and their ordinary world, highlighting " +
							"the current limitations and struggles in the field of energy production",
					},
					{
						Key:   "themeStated",
						Title: "The Importance of Responsible Innovation",
						Content: "A conversation among the scientists discussing the ethics and potential " +
							"consequences of tapping into extraordinary energy sources, foreshadowing the themes " +
							"of the story.",
					},
					{
						Key:   "setup",
						Title: "The World on the Brink of Change",
						Content: "Introduction to the main characters, their motivations, and the world they " +
							"live in, highlighting the team's pioneering work and the world's current energy crisis.",
					},
					{
						Key:   "catalyst",
						Title: "The Discovery of the Supernova Energy Source",
						Content: "A critical experiment or finding that reveals the possibility of harnessing " +
							"the energy of a nearby supernova, initiating the team's journey into the unknown.",
					},
					{
						Key:   "debate",
						Title: "The Dilemma of Power and Responsibility",
						Content: "The team grapples with the implications and potential fallout of their " +
							"discovery, voicing concerns and conflicting views on how to proceed, and setting the " +
							"stage for character arcs.",
					},
				},
				Lang:           models.LangEN,
				RegenerateKeys: []string{"openingImage", "debate", "setup"},
				UserID:         TestUser,
			},
		},
	}

	repository := daoai.NewRegenerateBeatsRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := lib.NewOpenaiContext(t.Context())

			beatsSheet, err := repository.RegenerateBeats(ctx, testCase.request)
			require.NoError(t, err)

			require.NotNil(t, beatsSheet)
			require.Len(t, beatsSheet, len(testCase.request.Plan.Beats))

			for i, beat := range beatsSheet {
				if lo.Contains(testCase.request.RegenerateKeys, beat.Key) {
					require.NotEqual(t, beat.Content, testCase.request.Beats[i].Content)
				} else {
					require.Equal(t, beat.Content, testCase.request.Beats[i].Content)
				}
			}

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

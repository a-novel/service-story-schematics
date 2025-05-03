package daoai_test

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/golm"
	groqbinding "github.com/a-novel-kit/golm/bindings/groq"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

func TestExpandBeat(t *testing.T) {
	testCases := []struct {
		name string

		request daoai.ExpandBeatRequest
	}{
		{
			name: "Success",

			request: daoai.ExpandBeatRequest{
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
				TargetKey: "setup",
				UserID:    TestUser,
			},
		},
	}

	binding := groqbinding.New(config.Groq.APIKey, config.Groq.Model)

	repository := daoai.NewExpandBeatRepository()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := golm.WithContext(t.Context(), binding)

			resp, err := repository.ExpandBeat(ctx, testCase.request)
			require.NoError(t, err)

			require.NotNil(t, resp)

			original, ok := lo.Find(testCase.request.Beats, func(item models.Beat) bool {
				return item.Key == testCase.request.TargetKey
			})
			require.True(t, ok)

			require.NotEqual(t, original.Content, resp.Content)

			CheckAgent(
				t,
				fmt.Sprintf(
					"Does the new beat expand on the original beat?\n\n"+
						"new sheet:\n\n%s\n\noriginal beat:\n\n%s",
					resp.Content, original,
				),
				fmt.Sprintf(
					"The new beat does not expand on the original beat.\n\n"+
						"new sheet:\n\n%s\n\noriginal beat:\n\n%s",
					resp.Content, original,
				),
			)
		})
	}
}

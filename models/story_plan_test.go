package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/story-schematics/models"
)

func TestBeatGetScenesCount(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		beat models.BeatDefinition

		expect string
	}{
		{
			name: "NoScenes",
		},
		{
			name: "OneSceneMinOnly",

			beat: models.BeatDefinition{
				MinScenes: 1,
			},

			expect: "(1 scene)",
		},
		{
			name: "OneSceneMaxOnly",

			beat: models.BeatDefinition{
				MaxScenes: 1,
			},

			expect: "(1 scene)",
		},
		{
			name: "OneSceneMinAndMax",

			beat: models.BeatDefinition{
				MinScenes: 1,
				MaxScenes: 1,
			},

			expect: "(1 scene)",
		},
		{
			name: "MultipleScenesMinOnly",

			beat: models.BeatDefinition{
				MinScenes: 3,
			},

			expect: "(3 scenes)",
		},
		{
			name: "MultipleScenesMaxOnly",

			beat: models.BeatDefinition{
				MaxScenes: 3,
			},

			expect: "(3 scenes)",
		},
		{
			name: "MultipleScenesMinAndMax",

			beat: models.BeatDefinition{
				MinScenes: 3,
				MaxScenes: 3,
			},

			expect: "(3 scenes)",
		},
		{
			name: "Range",

			beat: models.BeatDefinition{
				MinScenes: 2,
				MaxScenes: 5,
			},

			expect: "(2 to 5 scenes)",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, testCase.expect, testCase.beat.GetScenesCount())
		})
	}
}

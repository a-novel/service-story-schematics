package storyplanmodel_test

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

func TestScenesString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		scenes storyplanmodel.Scenes

		expect string
	}{
		{
			name: "Empty",

			scenes: storyplanmodel.Scenes{},

			expect: "any number of scenes",
		},
		{
			name: "ExactlyOne",

			scenes: storyplanmodel.Scenes{
				Exact: lo.ToPtr(1),
			},

			expect: "exactly 1 scene",
		},
		{
			name: "ExactlyMany",

			scenes: storyplanmodel.Scenes{
				Exact: lo.ToPtr(3),
			},

			expect: "exactly 3 scenes",
		},
		{
			name: "AtLeast",

			scenes: storyplanmodel.Scenes{
				Min: lo.ToPtr(2),
			},

			expect: "at least 2 scenes",
		},
		{
			name: "NoMoreThan",

			scenes: storyplanmodel.Scenes{
				Max: lo.ToPtr(5),
			},

			expect: "no more than 5 scenes",
		},
		{
			name: "Between",

			scenes: storyplanmodel.Scenes{
				Min: lo.ToPtr(2),
				Max: lo.ToPtr(4),
			},

			expect: "between 2 and 4 scenes",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := testCase.scenes.String()
			require.Equal(t, testCase.expect, actual)
		})
	}
}

func TestBeatOutputSchema(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		beat storyplanmodel.Beat

		expect any
	}{
		{
			name: "Default",

			beat: storyplanmodel.Beat{
				Name:      "Beat 1",
				Key:       "beat-1",
				KeyPoints: []string{"Key point 1", "Key point 2"},
				Purpose:   "Purpose of Beat 1",
				Scenes: storyplanmodel.Scenes{
					Exact: lo.ToPtr(3),
				},
			},

			expect: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"key", "content", "title"},
				"properties": map[string]any{
					"key": map[string]any{
						"const": "beat-1",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "A short title representing the beat.",
					},
					"content": map[string]any{
						"type": "string",
						"description": "A summary of the exactly 3 scenes that make up the 'Beat 1' beat." +
							"\nKey Points: " +
							"\n\t- Key point 1" +
							"\n\t- Key point 2" +
							"\nPurpose: Purpose of Beat 1",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := testCase.beat.OutputSchema()
			require.Equal(t, testCase.expect, actual)
		})
	}
}

func TestBeatString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		beat storyplanmodel.Beat

		expect string
	}{
		{
			name: "Default",

			beat: storyplanmodel.Beat{
				Name:      "Beat 1",
				Key:       "beat-1",
				KeyPoints: []string{"Key point 1", "Key point 2"},
				Purpose:   "Purpose of Beat 1",
				Scenes: storyplanmodel.Scenes{
					Exact: lo.ToPtr(3),
				},
			},

			expect: "Beat 1 (exactly 3 scenes)" +
				"\nKey Points: " +
				"\n\t- Key point 1" +
				"\n\t- Key point 2" +
				"\nPurpose: Purpose of Beat 1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := testCase.beat.String()
			require.Equal(t, testCase.expect, actual)
		})
	}
}

func TestPlanOutputSchema(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		plan *storyplanmodel.Plan

		expect any
	}{
		{
			name: "Default",

			plan: &storyplanmodel.Plan{
				Metadata: storyplanmodel.Metadata{},
				Beats: []storyplanmodel.Beat{
					{
						Name:      "Beat 1",
						Key:       "beat-1",
						KeyPoints: []string{"Key point 1", "Key point 2"},
						Purpose:   "Purpose of Beat 1",
						Scenes: storyplanmodel.Scenes{
							Exact: lo.ToPtr(3),
						},
					},
					{
						Name:      "Beat 2",
						Key:       "beat-2",
						KeyPoints: []string{"Key point 1"},
						Purpose:   "Purpose of Beat 2",
						Scenes: storyplanmodel.Scenes{
							Min: lo.ToPtr(2),
							Max: lo.ToPtr(4),
						},
					},
					{
						Name:      "Beat 3",
						Key:       "beat-3",
						KeyPoints: []string{"Key point 2"},
						Purpose:   "Purpose of Beat 3",
						Scenes: storyplanmodel.Scenes{
							Max: lo.ToPtr(4),
						},
					},
				},
			},

			expect: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"beats"},
				"properties": map[string]any{
					"beats": map[string]any{
						"type": "array",
						"description": "The beats that compose the story. " +
							"A beat is a unit of story structure that represents a specific moment or event in the " +
							"narrative.",
						"prefixItems": []any{
							map[string]any{
								"type":                 "object",
								"additionalProperties": false,
								"required":             []string{"key", "content", "title"},
								"properties": map[string]any{
									"key": map[string]any{
										"const": "beat-1",
									},
									"title": map[string]any{
										"type":        "string",
										"description": "A short title representing the beat.",
									},
									"content": map[string]any{
										"type": "string",
										"description": "A summary of the exactly 3 scenes that make up the 'Beat 1' beat." +
											"\nKey Points: " +
											"\n\t- Key point 1" +
											"\n\t- Key point 2" +
											"\nPurpose: Purpose of Beat 1",
									},
								},
							},
							map[string]any{
								"type":                 "object",
								"additionalProperties": false,
								"required":             []string{"key", "content", "title"},
								"properties": map[string]any{
									"key": map[string]any{
										"const": "beat-2",
									},
									"title": map[string]any{
										"type":        "string",
										"description": "A short title representing the beat.",
									},
									"content": map[string]any{
										"type": "string",
										"description": "A summary of the between 2 and 4 scenes that make up " +
											"the 'Beat 2' beat." +
											"\nKey Points: " +
											"\n\t- Key point 1" +
											"\nPurpose: Purpose of Beat 2",
									},
								},
							},
							map[string]any{
								"type":                 "object",
								"additionalProperties": false,
								"required":             []string{"key", "content", "title"},
								"properties": map[string]any{
									"key": map[string]any{
										"const": "beat-3",
									},
									"title": map[string]any{
										"type":        "string",
										"description": "A short title representing the beat.",
									},
									"content": map[string]any{
										"type": "string",
										"description": "A summary of the no more than 4 scenes that make up " +
											"the 'Beat 3' beat." +
											"\nKey Points: " +
											"\n\t- Key point 2" +
											"\nPurpose: Purpose of Beat 3",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual := testCase.plan.OutputSchema()
			require.Equal(t, testCase.expect, actual)
		})
	}
}

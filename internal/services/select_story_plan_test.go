package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

func TestSelectStoryPlan(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		request services.SelectStoryPlanRequest

		expect    *storyplanmodel.Plan
		expectErr error
	}{
		{
			name: "Success",

			request: services.SelectStoryPlanRequest{
				Lang: models.LangEN,
			},

			expect: storyplanmodel.SaveTheCat[models.LangEN],
		},
		{
			name: "NotFound",

			request: services.SelectStoryPlanRequest{
				Lang: models.Lang("unknown"),
			},

			expectErr: services.ErrStoryPlanNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			service := services.NewSelectStoryPlanService()

			resp, err := service.SelectStoryPlan(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)
		})
	}
}
